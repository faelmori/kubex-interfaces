package settings

//
//import (
//	//ida "github.com/faelmori/gospider/internal/middlewares"
//	//kbxApi "github.com/faelmori/kbxutils/api/config"
//	"github.com/faelmori/kbxutils/factory"
//	kbxSrv "github.com/faelmori/kbxutils/utils/interfaces"
//
//	ad "github.com/faelmori/gospider/api/data"
//	ar "github.com/faelmori/gospider/api/routes"
//	as "github.com/faelmori/gospider/api/schedule"
//	//ida "github.com/faelmori/gospider/internal/data/api"
//
//	//"github.com/gin-gonic/gin"
//	v "github.com/spf13/viper"
//
//	"fmt"
//	"path/filepath"
//	"time"
//)
//
//type configManager struct {
//	configPath string
//	keyPath    string
//	certPath   string
//	kbxConfig  kbxSrv.IConfigService
//	kbxFileSrv kbxSrv.FileSystemService
//}
//
//func NewOConfigManager(configPath string, keyPath string, certPath string) ad.IConfigManager {
//	flsrv := factory.NewFilesystemService(configPath)
//	cfgMgr := &configManager{
//		configPath: configPath,
//		keyPath:    keyPath,
//		certPath:   certPath,
//		kbxFileSrv: *flsrv,
//		kbxConfig:  factory.NewConfigService(configPath, keyPath, certPath),
//	}
//	if !cfgMgr.kbxConfig.IsConfigLoaded() {
//		if err := cfgMgr.kbxConfig.SetupConfig(); err != nil {
//			panic(fmt.Sprintf("Failed to setup config: %v", err))
//		}
//	}
//
//	return cfgMgr
//}
//
//func (cs *configManager) SetupConfig() error {
//	if !cs.kbxConfig.IsConfigLoaded() {
//		if setupErr := cs.kbxConfig.SetupConfig(); setupErr != nil {
//			return setupErr
//		}
//		v.SetConfigFile(cs.kbxFileSrv.GetConfigFilePath())
//		v.SetConfigType("json")
//		if err := v.ReadInConfig(); err != nil {
//			v.SetConfigFile(cs.kbxFileSrv.GetConfigFilePath())
//			if err2 := v.ReadInConfig(); err2 != nil {
//				return err2
//			}
//		}
//	}
//	return nil
//}
//
//func (cs *configManager) LoadConfig() error {
//	if loadErr := cs.kbxConfig.LoadConfig(); loadErr != nil {
//		return loadErr
//	}
//	cs.configPath = cs.kbxFileSrv.GetConfigFilePath()
//	if filepath.Base(cs.configPath) == "" {
//		cs.configPath = filepath.Join(cs.configPath, "config.json")
//	}
//	v.SetConfigFile(cs.configPath)
//	v.SetConfigType("json")
//	if err := v.ReadInConfig(); err != nil {
//		v.SetConfigFile(cs.kbxFileSrv.GetConfigFilePath())
//		if err2 := v.ReadInConfig(); err2 != nil {
//			return err2
//		}
//	}
//	return nil
//}
//
//func (cs *configManager) GetSetting(key string) (interface{}, error) {
//	if !v.IsSet(key) {
//		return nil, fmt.Errorf("setting %s not found", key)
//	}
//	return v.Get(key), nil
//}
//
//func (cs *configManager) UpdateCronScheduler(scheduler as.ICronScheduler) error {
//	var cronConfig map[string]interface{}
//	if err := v.UnmarshalKey("cron", &cronConfig); err != nil {
//		return err
//	}
//
//	for id, taskConfig := range cronConfig {
//		taskMap := taskConfig.(map[string]interface{})
//		action := func() error {
//			fmt.Println(taskMap["message"])
//			return nil
//		}
//		schedule, err := time.ParseDuration(taskMap["schedule"].(string))
//		if err != nil {
//			return err
//		}
//		recurring := taskMap["recurring"].(bool)
//		regTaskErr := scheduler.RegisterTask(id, action, schedule, recurring)
//		if regTaskErr != nil {
//			return regTaskErr
//		}
//	}
//
//	return nil
//}
//
//func (cs *configManager) UpdateRouter(rtr ar.IRouter) error {
//	var routesConfig map[string]interface{}
//	if err := v.UnmarshalKey("routes", &routesConfig); err != nil {
//		return err
//	}
//
//	//for path, routeConfig := range routesConfig {
//	//	routeMap := routeConfig.(map[string]interface{})
//	//	method := routeMap["method"].(string)
//	//	handler := func(c *gin.Context) {
//	//		c.JSON(200, gin.H{"message": routeMap["message"]})
//	//	}
//	//
//	//	route := NewR
//	//	WithMethod(method).
//	//		WithPath(path).
//	//		WithHandler(handler).
//	//		Build()
//	//
//	//	rtr.WithRoutes([]ar.IRoute{route})
//	//}
//
//	return nil
//}
//
//// Private methods (Internal)
//
//func (cs *configManager) GetConfigPath() string { return cs.configPath }
//func (cs *configManager) GetKeyPath() string    { return cs.keyPath }
//func (cs *configManager) GetCertPath() string   { return cs.certPath }
//
//// Private methods (External)
//
//func (cs *configManager) GetKbxConfigService() kbxSrv.IConfigService { return cs.kbxConfig }
