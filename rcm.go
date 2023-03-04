package rcm

func InitCenter(master bool, hostDomain string) {
	// master mode, serve as the domain 管理

}

func InitNode(totalTraffic int64, port int) *ServerNode {
	s := &ServerNode{}

	info := ServerInfo{
		TotalTraffic: totalTraffic,
		Port:         port,
		Unused:       TotalTraffic,
	}
	s.server = &info
	return s
}
