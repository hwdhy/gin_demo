package base

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// 项目初始化
func InitAppService() {
	//初始化配置
	initConfig()

	//初始化日志配置
	initLogs()

	// 从etcd中获取配置
	loadEtcdConf()

	//初始化mysql连接池
	initBaseMysql()

	//初始化session
	initSession()

	//初始化redis连接池
	initRedisPool()
}

// 初始化参数配置
func initConfig() {
	ConfFile = GetEnv("WEB_OVERSEAS_OFFICIAL_WEBSITE_ETCD", "/data/go/conf")                                                //etcd连接信息配置路径
	GConf.AppIdConfPath = GetEnv("WEB_APPID_CONF_PATH", "/data/go/conf")                                                    //APPID配置路径
	GConf.Env = GetEnv("WEB_OVERSEAS_OFFICIAL_WEBSITE_ENV", "dev")                                                          //环境变量, 和etcd中的保持一致
	GConf.Server.Host = GetEnv("WEB_OVERSEAS_OFFICIAL_WEBSITE_HOST", ":50051")                                              //http服务启动端口
	GConf.Server.LogRoot = GetEnv("WEB_OVERSEAS_OFFICIAL_WEBSITE_ROOT", "/data/wwwlogs/service/overseas_official_website/") //日志目录

	logrus.Infof("environment variables OVERSEAS_OFFICIAL_WEBSITE_ETCD value: %s, "+
		"WEB_OVERSEAS_OFFICIAL_WEBSITE_ENV value: %s, WEB_OVERSEAS_OFFICIAL_WEBSITE_HOST value: %s, "+
		"WEB_APPID_CONF_PATH value: %s", ConfFile, GConf.Env, GConf.Server.Host, GConf.AppIdConfPath)

	if ConfFile == "" || GConf.AppIdConfPath == "" || GConf.Env == "" || GConf.Server.Host == "" || GConf.Server.LogRoot == "" {
		logrus.Fatal("system environment variables value empty")
	}

	if ConfFile[len(ConfFile)-1:] != "/" {
		ConfFile = ConfFile + "/"
	}

	err := UnmarshalToml(ConfFile+"etcd-conf.toml", GConf)
	if err != nil {
		logrus.Fatal("fail to load conf file")
	}
	logrus.Infof("%+v", GConf)
}

// 日志配置初始化
func initLogs() {

	//初始化普通日志
	if GConf.Server.LogRoot[len(GConf.Server.LogRoot)-1:] != "/" {
		GConf.Server.LogRoot = GConf.Server.LogRoot + "/"
	}

	_, pathErr := os.Stat(GConf.Server.LogRoot)
	if pathErr != nil {
		makeErr := os.MkdirAll(GConf.Server.LogRoot, 0755)
		if makeErr != nil {
			logrus.Fatal("server log path make error")
		}
	}

	logrus.SetLevel(logrus.Level(GConf.Server.LogLevel))
	// 设置初始化全局日志
	MultipleLog = logrus.New()
	defaultLogPath := GConf.Server.LogRoot + "overseas_official_website"
	defaultLogFile, err := os.OpenFile(defaultLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.Fatalf("open file(%s) err:%v", defaultLogPath, err)
	}
	multiWriter := io.MultiWriter(defaultLogFile, os.Stdout)
	MultipleLog.SetOutput(multiWriter)

	MultipleLog.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})

	// 设置script日志
	scriptLogPath := GConf.Server.LogRoot + "script"
	scriptLogFile, err := os.OpenFile(scriptLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.Fatalf("open file(%s) err:%v", scriptLogPath, err)
	}

	ScriptLog = logrus.New()
	ScriptLog.SetOutput(scriptLogFile)
	ScriptLog.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})

	//输出到控制台+日志文件的日志
	MultipleLog.Info("Start overseas_official_website service...")
	MultipleLog.Infof("research host (%s)", GConf.Server.Host)
}

func GetEnv(path, defVal string) string {
	result := os.Getenv(path)
	if len(result) == 0 {
		return defVal
	}
	return result
}

func UnmarshalToml(file string, out interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return toml.Unmarshal(content, out)
}

// GetAllConfFiles 读取path目录下所有.json文件,解析到 configMap中去
func GetAllConfFiles(path string) (map[int64]AppIdConfig, error) {
	configMap := map[int64]AppIdConfig{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := file.Name()
		lowerName := strings.ToLower(name)
		if strings.Index(lowerName, ".json") == -1 || lowerName[len(lowerName)-5:] != ".json" {
			continue
		}

		filename := path + "/" + name
		appInfo := AppIdConfig{}
		content, err := os.ReadFile(filename)
		if err != nil {
			logrus.Errorf("read name: %s error, continue", name)
			continue
		}
		err = json.Unmarshal(content, &appInfo)
		if err != nil {
			logrus.Errorf("name: %s parse error, continue", name)
			continue
		}
		logrus.Infof("load name: %s to configMap success", name)
		configMap[appInfo.AppId] = appInfo
	}
	if len(configMap) < 1 {
		return nil, errors.New("configMap length = 0")
	}

	return configMap, nil
}
