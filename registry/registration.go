package registry

type Registration struct {
	ServiceName      ServiceName   // 服务名
	ServiceURL       string        // 服务URL
	RequiredServices []ServiceName // 依赖的服务列表
	ServiceUpdateURL string        // 通知可用服务的 URL
}

type ServiceName string

// 储存已有服务名
const (
	Log     ServiceName = "Log Service"
	Weather ServiceName = "Weather Service"
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
