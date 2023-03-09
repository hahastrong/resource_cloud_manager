package rcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ServerNode struct {
	mu     sync.Mutex
	server *ServerInfo
}

type ServerInfo struct {
	IP           string
	Port         int
	TotalTraffic int64
	Unused       int64
	Used         int64
	LastSeen     time.Time
}

type Response struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
}

func (s *ServerInfo) SendHeartBeat() {
	// 1. get the net traffic

	// 2. calculate the used traffic from first

	// 3. upload the traffic info to server

	dataByte, err := json.Marshal(s)
	if err != nil {
		return
	}

	res, err := http.Post(RegistryCenter, "application/json", bytes.NewReader(dataByte))
	if err != nil {
		return
	}

	var resByte []byte
	_, err = res.Body.Read(resByte)
	defer res.Body.Close()
	if err != nil {
		return
	}

	var response Response
	err = json.Unmarshal(resByte, &response)
	if err != nil {
		return
	}

	if response.Code != 0 {
		// 服务启动失败，因为接入服务太多了
		fmt.Println("too many servers")
		fmt.Println(response.Errmsg)
	}

	fmt.Printf("send heartbeat to center successful")

}

// collect traffic by update the used info
func (s *ServerNode) UpdateTrafficLoad(delta int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.server.Used += delta
	s.server.Unused -= delta
}

func (s *ServerNode) cronTaskToUploadTraffic() {
	c := time.Tick(time.Second * 10)
	for _ = range c {
		if s.server.Used > s.server.TotalTraffic {
			continue
		}
		s.server.SendHeartBeat()
	}
}
