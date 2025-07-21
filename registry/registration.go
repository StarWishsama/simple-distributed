package registry

type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

type ServiceName string

// 储存已有服务名
const (
	LogService ServiceName = "LogService"
)
