/* *
 * Test Dao
 */

package dao

import (
	"context"
	"io"

	"github.com/tonly18/xws/core/xerror"
)

type TestDao struct {
	ctx context.Context
	db  *dbConn
	rd  *redisPoolConn
}

func NewTestDao(ctx context.Context) *TestDao {
	return &TestDao{
		ctx,
		NewDBConn(ctx),
		NewRedis(ctx),
	}
}

func (d *TestDao) GetDB() *dbConn {
	return d.db
}

func (d *TestDao) GetRD() *redisPoolConn {
	return d.rd
}

func (d *TestDao) GetData(x int) (string, xerror.Error) {
	if x == 0 {
		//return "", xerror.NewXError("test-dao-error")
		return "", xerror.Wrap(io.EOF, "test-dao-error")
	}

	return "test dao", nil
}
