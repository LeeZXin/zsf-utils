package manage

import (
	"github.com/LeeZXin/zsf-utils/luautil"
	"github.com/LeeZXin/zsf-utils/maputil"
)

var (
	scriptFeatureConfigMap = maputil.NewConcurrentMap[string, *luautil.CachedScript](nil)
)

func RefreshScriptFeatureConfigMap(configMap map[string]*luautil.CachedScript) {
	scriptFeatureConfigMap.Clear()
	for k, v := range configMap {
		if v != nil {
			scriptFeatureConfigMap.Store(k, v)
		}
	}
}

func LoadScriptFeatureConfig(featureKey string) (*luautil.CachedScript, bool) {
	return scriptFeatureConfigMap.Load(featureKey)
}
