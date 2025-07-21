package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherInfo struct {
	City        string
	Temperature float32
	Description string
}

func RegisterHandlers() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			queryParams := r.URL.Query()
			city := queryParams.Get("city")
			if city == "" {
				http.NotFound(w, r)
			}

			for _, data := range mockWeatherDatas {
				if city == data.City {
					marshal, err := json.Marshal(data)

					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")

					if _, err = fmt.Fprint(w, string(marshal)); err != nil {
						return
					}
				}
			}

			w.WriteHeader(http.StatusNotFound)
			if _, err := fmt.Fprint(w, "Weather not found"); err != nil {
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
