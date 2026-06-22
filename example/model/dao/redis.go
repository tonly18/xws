package dao

import (
	"context"

	"github.com/extra-time-zone/xws/core/database"
	"github.com/redis/go-redis/v9"
)

type redisPoolConn struct {
	ctx context.Context
	rd  *redis.ClusterClient
}

func NewRedis(ctx context.Context) *redisPoolConn {
	return &redisPoolConn{
		ctx: ctx,
		rd:  database.GetRD(),
	}
}

func (d *redisPoolConn) GetRD() *redis.ClusterClient {
	return d.rd
}
