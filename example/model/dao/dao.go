package dao

import (
	"github.com/tonly18/xws/core/database"
	"github.com/tonly18/xws/example/pkg/global"
	"github.com/tonly18/xws/example/sconf"
)

func Init() {
	// mysql
	if err := initMysql(); err != nil {
		panic(err)
	}

	// redis
	if err := initRedis(); err != nil {
		panic(err)
	}
}

func initMysql() error {
	dbConf := make([]*database.MySqlConfig, 0, len(sconf.Config.MySql))
	for _, v := range sconf.Config.MySql {
		//master db
		mdb := makeDBConfig(v)
		//slave db
		sdb := make([]*database.MySqlConfig, 0, len(v.Slave))
		for _, s := range v.Slave {
			sdb = append(sdb, makeDBConfig(s))
		}

		mdb.Slave = sdb
		dbConf = append(dbConf, mdb)
	}

	if len(dbConf) > 0 {
		return database.InitDB(dbConf)
	}
	return nil
}

func initRedis() error {
	rdConf := makeRDConfig(sconf.Config.Redis)
	if global.RequiredTLS() {
		rdConf.TLS = true
		rdConf.ServerName = sconf.Config.Redis.ServerName
	}
	return database.InitRedis(rdConf)
}

func Close() {
	database.CloseDB()
	database.CloseRD()
}

func makeDBConfig(conf *sconf.MySqlConfig) *database.MySqlConfig {
	return &database.MySqlConfig{
		Role:         conf.Role,
		Host:         conf.Host,
		Port:         conf.Port,
		Dbname:       conf.Dbname,
		Username:     conf.Username,
		Password:     conf.Password,
		Charset:      conf.Charset,
		Collation:    conf.Collation,
		MaxIdleConns: conf.MaxIdleConns,
		MaxOpenConns: conf.MaxOpenConns,
		MaxLifetime:  conf.MaxLifetime,
		MaxIdleTime:  conf.MaxIdleTime,
		Slave:        nil,
	}
}

func makeRDConfig(conf *sconf.RedisConfig) *database.RedisClusterConfig {
	return &database.RedisClusterConfig{
		Host:         conf.Host,
		Username:     conf.Username,
		Password:     conf.Password,
		MaxIdleConns: conf.MaxIdleConns,
		MinIdleConns: conf.MinIdleConns,
		PoolSize:     conf.PoolSize,
		TLS:          false,
		ServerName:   "",
	}
}
