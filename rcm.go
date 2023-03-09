package rcm

func InitCenter() *Registry {
	// master mode, serve as the domain 管理
	r := NewRegistry()
	return r
}

func InitNode(totalTraffic int64, port int) *ServerNode {
	s := &ServerNode{}

	info := ServerInfo{
		TotalTraffic: totalTraffic,
		Port:         port,
		Unused:       totalTraffic,
	}
	s.server = &info

	// cron task to upload traffic info
	go s.cronTaskToUploadTraffic()

	return s
}
