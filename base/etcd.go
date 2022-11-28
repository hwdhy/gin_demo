package base

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func loadEtcdConf() {
	//获取此路径下的所有配置
	keyPrefix := fmt.Sprintf("%s/overseas_official_website/%s", GConf.ETCD.RootKey, GConf.Env)

	//具体的路径 key解析到相应的结构体内
	etcdConfigMap := map[string]interface{}{
		fmt.Sprintf(EtcdCommonBase, GConf.ETCD.RootKey, GConf.Env):         &GConf.CommonBase,
		fmt.Sprintf(EtcdCommonHttp, GConf.ETCD.RootKey, GConf.Env):         &GConf.CommonHttp,
		fmt.Sprintf(EtcdCommonMysqlTimeout, GConf.ETCD.RootKey, GConf.Env): &GConf.CommonMysqlTimeout,
		fmt.Sprintf(EtcdCommonRedisTimeout, GConf.ETCD.RootKey, GConf.Env): &GConf.CommonRedisTimeout,
		fmt.Sprintf(EtcdCommonRefresh, GConf.ETCD.RootKey, GConf.Env):      &GConf.CommonRefresh,
		fmt.Sprintf(EtcdMysqlMaster, GConf.ETCD.RootKey, GConf.Env):        &GConf.MysqlMaster,
		fmt.Sprintf(EtcdMysqlSlave, GConf.ETCD.RootKey, GConf.Env):         &GConf.MysqlSlave,
		fmt.Sprintf(EtcdRedisConfig, GConf.ETCD.RootKey, GConf.Env):        &GConf.RedisConfig,
	}
	GetEtcdConfig(GConf.ETCD, keyPrefix, etcdConfigMap, MultipleLog)
}

// GetEtcdConfig 连接ETCD, 根据路径前缀获取所有KV, 根据全路径匹配将配置解析到相应的结构体内
func GetEtcdConfig(etcdConfig EtcdConfig, keyPrefix string, configMap map[string]interface{}, multipleLog *logrus.Logger) {
	//连接ETCD
	dialTimeout := time.Duration(etcdConfig.DialTimeout) * time.Second
	clientConfig := clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   etcdConfig.Endpoints,
		Username:    etcdConfig.Username,
		Password:    etcdConfig.Password,
	}
	//连接etcd
	etcdCli, err := clientv3.New(clientConfig)
	if err != nil {
		multipleLog.Fatalf("init etcd failed, err: %v", err)
	}
	multipleLog.Info("ETCD connect success")

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	//判断etcd状态
	_, err = etcdCli.Status(ctx, clientConfig.Endpoints[0])
	if err != nil {
		multipleLog.Fatalf("error checking etcd endpoints 0 status, err: %v", err)
	}

	//取值
	value, err := etcdCli.Get(ctx, keyPrefix, clientv3.WithPrefix())
	defer cancel()
	etcdKvs := map[string][]byte{}
	if err != nil {
		multipleLog.Fatalf("get etcd key: %s failed, err: %v", keyPrefix, err)
	}
	multipleLog.Infof("ETCD get value by key prefix: %s success", keyPrefix)

	if value.Count == 0 {
		multipleLog.Fatalf("get etcd key: %s failed, count = 0", keyPrefix)
	}
	for _, resp := range value.Kvs {
		etcdKvs[string(resp.Key)] = resp.Value
	}

	//赋值
	for k, v := range configMap {
		if _, ok := etcdKvs[k]; !ok {
			multipleLog.Fatalf("get etcd key: %s, don't exits, stop, error: %v", k, err)
		}
		err = json.Unmarshal(etcdKvs[k], v)
		if err != nil {
			multipleLog.Fatalf("get etcd key: %s, error: %v, stop", k, err)
		}
		validate := validator.New()
		err = validate.Var(v, "required")
		if err != nil {
			multipleLog.Fatalf("get etcd key: %s, value: %+v, error: %v, validator failed, stop", k, v, err)
		}
		multipleLog.Infof("key: %s validator success", k)
	}

	multipleLog.Info("ETCD load data to config success")
}
