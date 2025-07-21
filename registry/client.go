package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-distributed/util"
)

var registryUrl string

func init() {
	v, err := util.InitViper("regservice")

	if err != nil {
		panic("Failed to initialize configuration: " + err.Error())
	}

	registryUrl = "http://" + v.GetString("server.host") + ":" + v.GetString("server.port") + "/services"
}

// RegisterService 向注册中心注册服务
func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)

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
