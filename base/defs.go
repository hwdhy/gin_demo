package base

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

// ConfFile 配置文件路径
var ConfFile string

// 配置
var GConf = &GameConf{}

// 主库Db
var MasterDb *gorm.DB

// 从库Db
var SlaveDb *gorm.DB

// Redis连接池
var RedisPool *redis.Pool

// 定时任务的日志文件
var ScriptLog *logrus.Logger

// 输出到控制台+日志文件
var MultipleLog *logrus.Logger

// 内存中存储信息
var MemoryStoreInfo sync.Map

const (
	// EtcdCommonBase 基本信息
	EtcdCommonBase = "%s/overseas_official_website/%s/common/base"

	// EtcdCommonHttp http超时配置
	EtcdCommonHttp = "%s/overseas_official_website/%s/common/http"

	// EtcdCommonMysqlTimeout mysql超时配置
	EtcdCommonMysqlTimeout = "%s/overseas_official_website/%s/common/mysql/timeout"

	// EtcdCommonRedisTimeout redis超时配置
	EtcdCommonRedisTimeout = "%s/overseas_official_website/%s/common/redis/timeout"

	// EtcdCommonRefresh 定时任务配置
	EtcdCommonRefresh = "%s/overseas_official_website/%s/common/refresh"

	// EtcdMysqlMaster 主库配置
	EtcdMysqlMaster = "%s/overseas_official_website/%s/mysql/master"

	// EtcdMysqlSlave 从库配置
	EtcdMysqlSlave = "%s/overseas_official_website/%s/mysql/slave"

	// EtcdRedisConfig redis连接配置
	EtcdRedisConfig = "%s/overseas_official_website/%s/redis/config"
)

// EtcdConfig 配置结构体
type EtcdConfig struct {
	DialTimeout int64
	Endpoints   []string
	Username    string
	Password    string
	RootKey     string //根路径
}

// 服务器通用配置内容
type ServerConf struct {
	ServiceConf

	LogRoot      string
	LogLevel     int8
	MaxProcs     int
	DataLogsPath string
}

type ServiceConf struct {
	// 例如，同一台机器可能为多个项目部署同一个登录服务的多个实例，名字可以填为项目id
	Name string
	// 运维为每个服务指派的id。某些场景，一台物理机可能部署同一个服务的多个实例，所以需要id区分（例如同一台机器，部署4个登录服务）
	ID int64
	//http Listen host
	Host string
}

// 配置信息
type GameConf struct {
	Server        ServerConf
	ETCD          EtcdConfig
	Env           string
	AppIdConfPath string

	CommonBase         CommonBase
	CommonHttp         CommonHttp
	CommonMysqlTimeout CommonMysqlTimeout
	CommonRedisTimeout CommonRedisTimeout
	CommonRefresh      CommonRefresh
	//mysql
	MysqlMaster MysqlConfig
	MysqlSlave  MysqlConfig
	//redis
	RedisConfig RedisConfig
}

type CommonBase struct {
	HostileSaveCheckDuration int   `validate:"required"` //单位时间,秒
	HostileSaveCheckMaxTimes int   `validate:"required"` //单位时间内,恶意提交次数限制
	SessionExpired           int64 `validate:"required"` //Session过期时间
	CallRestfulTimeout       int   `validate:"required"`
}

type CommonHttp struct {
	ReadTimeout  int `validate:"required"`
	WriteTimeout int `validate:"required"`
	IdleTimeout  int `validate:"required"`
}

type CommonMysqlTimeout struct {
	MysqlTimeout      string `validate:"required"`
	MysqlReadTimeout  string `validate:"required"`
	MysqlWriteTimeout string `validate:"required"`
}

type CommonRedisTimeout struct {
	//redis dial timeout
	RedisDialConnectTimeout int64 `validate:"required"`
	RedisDialReadTimeout    int64 `validate:"required"`
	RedisDialWriteTimeout   int64 `validate:"required"`
}

type CommonRefresh struct {
	AppKeyRefreshTime int `validate:"required"`
}

// 阿里短信服务配置
type AliSmsServiceConfig struct {
	RegionId  string
	AccessId  string
	SecretKey string
}

// 阿里短信模板配置
type AliSmsTemplateConfig struct {
	TemplateCode string
	SignName     string
}

// 邮件配置信息
type MailConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Charset  string
}

// 邮件模板
type MailTpl struct {
	Type    string
	LangId  string
	Title   string
	Content string
}

// redis配置信息结构
type RedisConfig struct {
	Host        string
	Port        int
	Pass        string
	Db          int
	MaxActive   int
	MaxIdle     int
	IdleTimeout int
}

// mysql配置信息结构
type MysqlConfig struct {
	Host        string
	User        string
	Port        int
	Pass        string
	MaxConn     int
	MaxIdle     int
	MaxLifetime int
	DbName      string
	Charset     string
}

// 输出信息
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// AppIdConfig SDK/服务器/客户端调用接口签名时所用的app_id, secret key等信息
type AppIdConfig struct {
	GameId     int    `json:"game_id" validate:"required"`
	AppId      int64  `json:"app_id" validate:"required"`
	SecretKey  string `json:"secret_key" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Type       int    `json:"type" validate:"required"` //1 为SDK, 2 为服务器， 3 为客户端
	Enabled    int    `json:"enabled" validate:"required"`
	PlatformId int64  `json:"platform_id" validate:"required"`
}
