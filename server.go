package rcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ServerNode struct {
	mu     sync.Mutex
	server *ServerInfo
}

type ServerInfo struct {
	IP            string
	Port          int
	TotalTraffic  int64
	Unused        int64
	Used          int64
	LastSeen      time.Time
	ResourceReady bool
	LastUpdate    time.Time
}

type Response struct {
	Code         int    `json:"code"`
	Errmsg       string `json:"errmsg"`
	ResourceLink string `json:"resource_link"`
}

func (s *ServerInfo) SendHeartBeat() error {
	// 1. get the net traffic

	// 2. calculate the used traffic from first

	// 3. upload the traffic info to server

	dataByte, err := json.Marshal(s)
	if err != nil {
		return err
	}

	res, err := http.Post(RegistryCenter, "application/json", bytes.NewReader(dataByte))
	if err != nil {
		return err
	}

	var resByte []byte
	resByte, err = io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(resByte, &response)
	if err != nil {
		return err
	}

	if response.Code != 0 {
		// 服务启动失败，因为接入服务太多了
		fmt.Println("too many servers")
		fmt.Println(response.Errmsg)
		return errors.New(response.Errmsg)
	}

	fmt.Printf("send heartbeat to center successful")

	// download the resource
	s.DownloadResource(response.ResourceLink)

	return nil

}

func (s *ServerInfo) DownloadResource(path string) {
	if len(path) == 0 {
		return
	}

	// 一旦收到注册中心下发的数据，就下载资源
	s.ResourceReady = false

	err := DownloadDirFiles(path)
	if err != nil {
		return
	}

	s.ResourceReady = true
	s.LastUpdate    = time.Now()

	// delete the expired resource
	DeleteExpiredResource()
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
		_ = s.server.SendHeartBeat()
	}
}
