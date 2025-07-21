package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"simple-distributed/util"
	"slices"
	"sync"
)

var registryUrl string

func init() {
	v, err := util.InitViper("regservice")

	if err != nil {
		panic("Failed to initialize configuration: " + err.Error())
	}

	registryUrl = "http://" + v.GetString("server.host") + ":" + v.GetString("server.port") + "/services"
}

// serviceUpdateHandler 服务更新处理器
type serviceUpdateHandler struct{}

func (h *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(r.Body)

	var p patch
	if err := dec.Decode(&p); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadGateway)
		return
	}

	if err := prov.Update(p); err != nil {
		log.Println(err)
	}
}

// RegisterService 向注册中心注册服务
func RegisterService(r Registration) error {
	srvUpdateURL, err := url.Parse(r.ServiceUpdateURL)

	if err != nil {
		return err
	}

	http.Handle(srvUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)

	if err != nil {
		return err
	}

	resp, err := http.Post(registryUrl, "application/json", buf)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service %s responded with %s", r.ServiceName, resp.StatusCode)
	}

	return nil
}

// UnregisterService 从注册中心注销服务
func UnregisterService(url string) error {
	req, err := http.NewRequest("DELETE", registryUrl, bytes.NewBufferString(url))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/plain")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unregister service %s responded with %s", url, resp.Status)
	}

	return nil
}

// providers 用于存储当前服务的提供者信息
type providers struct {
	services map[ServiceName][]string // 提供者的 URL 列表, 一个服务可能有多个 URL 可供使用
	mutex    *sync.RWMutex
}

var prov = providers{services: make(map[ServiceName][]string), mutex: new(sync.RWMutex)}

// Update 更新服务提供者信息
func (p *providers) Update(pat patch) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, pe := range pat.Added {
		if _, exists := p.services[pe.Name]; !exists {
			p.services[pe.Name] = make([]string, 0)
		}

		p.services[pe.Name] = append(p.services[pe.Name], pe.URL)
	}

	for _, pe := range pat.Removed {
		if srvs, exists := p.services[pe.Name]; exists {
			for i, srv := range srvs {
				if srv == pe.URL {
					p.services[pe.Name] = slices.Delete(srvs, i, i+1)
				}
			}
		}
	}

	return nil
}

// get 根据服务名获取服务提供者的 URL 列表
func (p *providers) get(name ServiceName) ([]string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if urls, exists := p.services[name]; exists {
		return urls, nil
	}

	return nil, fmt.Errorf("service %s not found", name)
}

func GetProvider(name ServiceName) (string, error) {
	urls, err := prov.get(name)

	if err != nil {
		return "", err
	}

	if len(urls) == 0 {
		return "", fmt.Errorf("no providers available for service %s", name)
	}

	randIdx := rand.Int() % len(urls) // FIXME: 可以用负载均衡的方式，而非随机获取

	return urls[randIdx], nil
}
