package dao

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/tonly18/xws/core/database"
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
