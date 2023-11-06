package manage

import (
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"github.com/LeeZXin/zsf-utils/luautil"
)

var (
	scriptFeatureConfigMap = hashmap.NewConcurrentHashMap[string, *luautil.CachedScript]()
)

func RefreshScriptFeatureConfigMap(configMap map[string]*luautil.CachedScript) {
	scriptFeatureConfigMap.Clear()
	for k, v := range configMap {
		if v != nil {
			scriptFeatureConfigMap.Put(k, v)
		}
	}
}

func LoadScriptFeatureConfig(featureKey string) (*luautil.CachedScript, bool) {
	return scriptFeatureConfigMap.Get(featureKey)
}
