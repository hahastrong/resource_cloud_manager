package rcm

import (
	"fmt"
	"sync"
	"time"
)

type Registry struct {
	mu      sync.Mutex
	servers map[string]*ServerInfo
}

func NewRegistry() *Registry {
	return &Registry{
		servers: make(map[string]*ServerInfo),
	}
}

func (s *Registry) RegisterOrUpdateInfo(server ServerInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, ok := s.servers[server.IP]
	if !ok {
		info = &server
		s.servers[server.IP] = info
	}
	info.Used = server.Used
	// 更新已用流量

	info.LastSeen = time.Now()
}

func (s *Registry) Unregister(serverIP string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.servers, serverIP)
}

func (s *Registry) ListServers() []ServerInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	var servers []ServerInfo
	for _, info := range s.servers {
		servers = append(servers, *info)
	}

	return servers
}

func (s *Registry) SelectServer() (*ServerInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var selected *ServerInfo
	// random select
	for _, info := range s.servers {
		if info.Used < info.TotalTraffic {
			selected = info
			break
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no servers available")
	}

	return selected, nil
}

func (s *Registry) CleanInactiveServer() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var deletedKeys []string

	for key, info := range s.servers {
		if time.Now().Sub(info.LastSeen) > time.Minute * 30 {
			deletedKeys = append(deletedKeys, key)
		}
	}

	for _, key := range deletedKeys {
		delete(s.servers, key)
	}
}
