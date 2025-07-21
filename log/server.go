package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

// fileLog 代表以文件形式存储的日志
type fileLog string

// Write 写入日志到文件
// implements the io.Writer interface
func (fl fileLog) Write(data []byte) (int, error) {
	file, err := os.OpenFile(string(fl), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(data)
}

// Run 初始化日志记录器
func Run(dest string) {
	log = stlog.New(fileLog(dest), "go: ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)

			if err != nil || len(msg) == 0 {
				http.Error(w, "Invalid log message", http.StatusBadRequest)
				return
			}

			write(string(msg))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(msg string) {
	log.Printf("%v\n", msg)
}
