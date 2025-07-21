package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type registry struct {
	registrations []Registration // 服务注册列表
	mutex         *sync.RWMutex  // 读写互斥锁，确保并发安全
}

// Register 将服务注册到注册中心
func (r *registry) Register(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registrations = append(r.registrations, reg)

	if err := r.sendRequiredService(reg); err != nil {
		return err
	}

	return nil
}

// Unregister 从注册中心注销服务
func (r *registry) Unregister(url string) error {
	for i, reg := range r.registrations {
		if reg.ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()

			return nil
		}
	}

	return fmt.Errorf("the service of %v was not found", url)
}

// sendRequiredService 发送该服务所需的服务信息
func (r *registry) sendRequiredService(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 创建一个补丁对象，用于存储服务依赖的其他服务
	var p patch

	for _, srvReg := range r.registrations {
		for _, need := range reg.RequiredServices {
			if srvReg.ServiceName == need {
				p.Added = append(p.Added, patchEntry{
					Name: srvReg.ServiceName,
					URL:  srvReg.ServiceURL,
				})
			}
		}
	}

	if err := r.sendPatch(p, reg.ServiceUpdateURL); err != nil {
		return err
	}

	return nil
}

// sendPatch 发送服务中心服务变动的补丁信息
func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)

	if err != nil {
		return err
	}

	if _, err = http.Post(url, "application/json", bytes.NewBuffer(d)); err != nil {
		return err
	}

	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			http.Error(w, "Invalid registration data", http.StatusBadRequest)
			return
		}

		log.Printf("Registering service: %s at %s\n", r.ServiceName, r.ServiceURL)

		err = reg.Register(r)

		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to register service", http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		url := string(payload)

		log.Printf("Unregistering service at %s\n", url)
		err = reg.Unregister(url)

		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to unregister service", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
