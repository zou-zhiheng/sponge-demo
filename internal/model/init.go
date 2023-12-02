// Package model is the initial database driver and define the data structures corresponding to the tables.
package model

import (
	"errors"
	"fmt"
	"github.com/zhuxiujia/GoMybatis"
	"os"
	"sync"
	"time"

	"user/internal/config"

	"github.com/zhufuyi/sponge/pkg/goredis"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// ErrCacheNotFound No hit cache
	ErrCacheNotFound = redis.Nil

	// ErrRecordNotFound no records found
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	db    *gorm.DB
	once1 sync.Once

	redisCli *redis.Client
	once2    sync.Once

	cacheType *CacheType
	once3     sync.Once

	mybatisEngine *GoMybatis.GoMybatisEngine
	once4         sync.Once

	userMapper *UserMapper
	once5      sync.Once
)

// InitMysql connect mysql
func InitMysql() {
	opts := []mysql.Option{
		mysql.WithMaxIdleConns(config.Get().Mysql.MaxIdleConns),
		mysql.WithMaxOpenConns(config.Get().Mysql.MaxOpenConns),
		mysql.WithConnMaxLifetime(time.Duration(config.Get().Mysql.ConnMaxLifetime) * time.Minute),
	}
	if config.Get().Mysql.EnableLog {
		opts = append(opts,
			mysql.WithLogging(logger.Get()),
			mysql.WithLogRequestIDKey("request_id"),
		)
	}

	if config.Get().App.EnableTrace {
		opts = append(opts, mysql.WithEnableTrace())
	}

	// setting mysql slave and master dsn addresses,
	// if there is no read/write separation, you can comment out the following piece of code
	opts = append(opts, mysql.WithRWSeparation(
		config.Get().Mysql.SlavesDsn,
		config.Get().Mysql.MastersDsn...,
	))

	// add custom gorm plugin
	//opts = append(opts, mysql.WithGormPlugin(yourPlugin))

	var err error
	db, err = mysql.Init(config.Get().Mysql.Dsn, opts...)
	if err != nil {
		panic("mysql.Init error: " + err.Error())
	}
}

// GetDB get db
func GetDB() *gorm.DB {
	if db == nil {
		once1.Do(func() {
			InitMysql()
		})
	}

	return db
}

// CloseMysql close mysql
func CloseMysql() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}

// ------------------------------------------------------------------------------------------

// CacheType cache type
type CacheType struct {
	CType string        // cache type  memory or redis
	Rdb   *redis.Client // if CType=redis, Rdb cannot be empty
}

// InitCache initial cache
func InitCache(cType string) {
	cacheType = &CacheType{
		CType: cType,
	}

	if cType == "redis" {
		cacheType.Rdb = GetRedisCli()
	}
}

// GetCacheType get cacheType
func GetCacheType() *CacheType {
	if cacheType == nil {
		once3.Do(func() {
			InitCache(config.Get().App.CacheType)
		})
	}

	return cacheType
}

// InitRedis connect redis
func InitRedis() {
	opts := []goredis.Option{
		goredis.WithDialTimeout(time.Duration(config.Get().Redis.DialTimeout) * time.Second),
		goredis.WithReadTimeout(time.Duration(config.Get().Redis.ReadTimeout) * time.Second),
		goredis.WithWriteTimeout(time.Duration(config.Get().Redis.WriteTimeout) * time.Second),
	}
	if config.Get().App.EnableTrace {
		opts = append(opts, goredis.WithEnableTrace())
	}

	var err error
	redisCli, err = goredis.Init(config.Get().Redis.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedisCli get redis client
func GetRedisCli() *redis.Client {
	if redisCli == nil {
		once2.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis close redis
func CloseRedis() error {
	if redisCli == nil {
		return nil
	}

	err := redisCli.Close()
	if err != nil && err.Error() != redis.ErrClosed.Error() {
		return err
	}

	return nil
}

func InitGoMybatis() {
	tmp := GoMybatis.GoMybatisEngine{}.New()
	mybatisEngine = &tmp

	mybatisEngine.TemplateDecoder().SetPrintElement(false)
	mybatisEngine.SetPrintWarning(false)

	_, err := mybatisEngine.Open("mysql", config.Get().Mysql.Dsn)
	if err != nil {
		panic(err)
	}

	//动态数据源路由(可选)
	/**
	engine.Open("mysql", MysqlUri)//添加第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
	var router = GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
		//根据包名路由指向数据源
		if strings.Contains(mapperName, "example.") {
			var url = MysqlUri//第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
			fmt.Println(url)
			return &url
		}
		return nil
	})
	engine.SetDataSourceRouter(&router)
	**/

	//自定义日志实现(可选)
	//engine.SetLogEnable(true)
	//engine.SetLog(&GoMybatis.LogStandard{
	//	PrintlnFunc: func(messages ...string) {
	//		println("log>> ", fmt.Sprint(messages))
	//	},
	//})
}

func GetMybatisEngine() *GoMybatis.GoMybatisEngine {
	if mybatisEngine == nil {
		once4.Do(func() {
			InitGoMybatis()
		},
		)
	}

	return mybatisEngine
}

func CloseMybatisEngine() error {

	if mybatisEngine == nil {
		return errors.New("MybatisEngine has been closed")
	}

	mybatisEngine = nil

	return nil
}

func InitUserMapper() {

	engine := GetMybatisEngine()
	wd, _ := os.Getwd()
	fmt.Println("wd", wd)
	//bytes, err := os.ReadFile("./internal/model/mapper/UserMapper.xml")
	bytes, err := os.ReadFile("../model/mapper/UserMapper.xml")
	if err != nil {
		panic(err)
	}

	userMapper = &UserMapper{}
	engine.WriteMapperPtr(userMapper, bytes)

}

func GetUserMapper() *UserMapper {
	if userMapper == nil {
		once5.Do(func() {
			InitUserMapper()
		})
	}

	return userMapper
}

func CloseUserMapper() error {

	if userMapper == nil {
		return errors.New("UserMapper has been closed")
	}

	userMapper = nil

	return nil
}
