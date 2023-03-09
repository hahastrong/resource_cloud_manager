package rcm

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := InitNode(123000, 7999)

	tt := time.Tick(time.Second * 3)

	for range tt {
		s.UpdateTrafficLoad(123)
	}

	time.Sleep(time.Minute * 10)
}
