package ziface

import "golang.org/x/time/rate"

type IRater interface {
	GetLimiter() *rate.Limiter
}
