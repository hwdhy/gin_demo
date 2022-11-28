package refresh

import (
	"fmt"
	"overseas-official-website/base"
	"time"
)

func AppKeyRefresh() {
	loadApiAppKey(true)

	for range time.Tick(time.Second * time.Duration(base.GConf.CommonRefresh.AppKeyRefreshTime)) {
		loadApiAppKey(false)
	}
}

// 定时刷新app_id和secret_key
func loadApiAppKey(stopProcess bool) {
	appIdMap, err := base.GetAllConfFiles(base.GConf.AppIdConfPath)
	if err != nil {
		//如果时启动时加载, 有错误则停止
		if stopProcess {
			base.ScriptLog.Fatalf("LoadApiAppKey, load appid config error: %s", err.Error())
			return
		}
		base.ScriptLog.Errorf("LoadApiAppKey, load appid config error: %s", err.Error())
		return
	}
	for appId, appInfo := range appIdMap {
		if appInfo.Type != 3 {
			continue
		}
		key := fmt.Sprintf("app_id_%d", appId)
		base.MemoryStoreInfo.Store(key, appInfo)
		base.ScriptLog.Infof("LoadApiAppKey, success load into memory: %s", key)
	}
}
