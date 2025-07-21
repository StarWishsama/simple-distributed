package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type registry struct {
	registrations []Registration // 服务注册列表
	mutex         *sync.Mutex    // 互斥锁，确保并发安全
}

func (r *registry) Register(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registrations = append(r.registrations, reg)
	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
