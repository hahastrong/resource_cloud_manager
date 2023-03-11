package rcm

import (
	"fmt"
	"math/rand"
)

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

	// set the random delta time
	RandomDeltaUpdate = int64(rand.Int()) % 600

	err := s.server.SendHeartBeat()
	if err != nil {
		fmt.Println("failed to registry the node")
		panic(err)
	}

	// cron task to upload traffic info
	go s.cronTaskToUploadTraffic()

	return s
}
