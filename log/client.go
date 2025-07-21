package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"
	"simple-distributed/registry"
)

type clientLogger struct {
	url string
}

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] -", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

func (c *clientLogger) Write(data []byte) (n int, err error) {
	b := bytes.NewBuffer(data)

	resp, err := http.Post(c.url+"/log", "text/plain", b)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to log message, log service responded with %s", resp.Status)
	}

	return len(data), nil
}
