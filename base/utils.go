package base

import (
	"fmt"
	sessionRedis "github.com/go-session/redis"
	"github.com/go-session/session"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 初始化Session redis
func initSession() {
	config := GConf.RedisConfig
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	session.InitManager(
		session.SetCookieName("hrg_user_research"),
		session.SetExpired(GConf.CommonBase.SessionExpired),
		session.SetStore(sessionRedis.NewRedisStore(&sessionRedis.Options{
			Addr:     addr,
			Password: config.Pass,
			DB:       config.Db,
		})),
	)
}

func initBaseMysql() {
	MasterDb = connMysql(GConf.MysqlMaster)
	if MasterDb != nil {
		MultipleLog.Info("Mysql Master connect success")
	}
	SlaveDb = connMysql(GConf.MysqlSlave)
	if SlaveDb != nil {
		MultipleLog.Info("Mysql Slave connect success")
	}
}

// 连接数据库
func connMysql(config MysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true",
		config.User, config.Pass, config.Host, config.Port, config.DbName, GConf.CommonMysqlTimeout.MysqlTimeout,
		GConf.CommonMysqlTimeout.MysqlReadTimeout, GConf.CommonMysqlTimeout.MysqlWriteTimeout)

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		MultipleLog.Fatal(fmt.Sprintf("dsn: %s Open fail, error: %s", dsn, err.Error()))
	}
	db, _ := mysqlDb.DB()

	db.SetMaxOpenConns(config.MaxConn)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	err = db.Ping()
	if err != nil {
		MultipleLog.Fatal(fmt.Sprintf("dsn: %s Ping fail, error: %s", dsn, err.Error()))
		return nil
	}
	return mysqlDb
}
