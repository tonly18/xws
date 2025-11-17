package dao

import (
	"context"

	"github.com/tonly18/xws/core/database"
	"gorm.io/gorm"
)

type dbConn struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDBConn(ctx context.Context, roles ...string) *dbConn {
	role := database.RoleDefault
	if len(roles) > 0 {
		role = roles[0]
	}

	return &dbConn{
		ctx: ctx,
		db:  database.GetDB(role),
	}
}

func (d *dbConn) GetDB() *gorm.DB {
	return d.db
}
