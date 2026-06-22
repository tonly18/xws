package zutils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/extra-time-zone/xws/core/zconf"
)

type Signal struct {
	sigChan chan os.Signal
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSignal() *Signal {
	s := &Signal{
		sigChan: make(chan os.Signal),
		ctx:     nil,
		cancel:  nil,
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *Signal) notify() {
	signal.Notify(s.sigChan, syscall.SIGINT, syscall.SIGTERM)
}

func (s *Signal) Waiter() error {
	s.notify()

	for sig := range s.sigChan {
		switch sig {
		case syscall.SIGINT:
			fmt.Printf("Received control signal int: %+v\n", sig)
			return nil
		case syscall.SIGTERM:
			fmt.Printf("Received control signal term: %+v\n", sig)
			return nil
		default:
			fmt.Printf("Received control signal: %+v\n", sig)
		}
	}

	return nil
}

func (s *Signal) Cannel(duration ...time.Duration) {
	// 服务退出
	zconf.ServiceExit.Store(true)

	// cancel
	s.cancel()

	// 等待业务处理完成
	if len(duration) > 0 {
		time.Sleep(duration[0])
	}
}

func (s *Signal) Context() context.Context {
	return s.ctx
}
