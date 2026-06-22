package znet

import (
	"github.com/extra-time-zone/xws/core/ziface"
	"golang.org/x/time/rate"
)

const normalRate = 5 //平时速率
const maxRate = 50   //瞬间最大速率

type Rater struct {
	limiter *rate.Limiter
}

func NewRater() ziface.IRater {
	return &Rater{
		limiter: rate.NewLimiter(normalRate, maxRate),
	}
}

func (r *Rater) GetLimiter() *rate.Limiter {
	return r.limiter
}
