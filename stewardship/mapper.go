package stewardship

//import (
//	"errors"
//	"github.com/faelmori/kubex-interfaces/module"
//	l "github.com/faelmori/logz"
//	"os"
//	"path/filepath"
//	"strings"
//)
//
//// Mapper is a struct that provides methods to discover and manage modules configurations.
//type Mapper struct {
//	BasePath string
//}
//
//// NewMapper creates a new Mapper instance with the specified base path.
//func NewMapper(basePath string) *Mapper {
//	return &Mapper{BasePath: basePath}
//}
//
//// DiscoverModules walks the directory tree starting from BasePath and returns a list of module configurations found.
//func (m *Mapper) DiscoverModules() (map[string]module.KubexModule, error) {
//	var modulesConfig = make(map[string]module.KubexModule)
//
//	err := filepath.Walk(m.BasePath, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			l.ErrorCtx("ErrorCtx walking the path", map[string]interface{}{
//				"context":  "kbxutils",
//				"action":   "discoverModules",
//				"error":    err.Error(),
//				"path":     path,
//				"showData": true,
//			})
//			return err
//		}
//		if !info.IsDir() && strings.Compare(strings.ToLower(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))), "config") == 0 {
//			// Check if the file is a module configuration and if is already in the cache
//			if _, found := modulesConfig[info.Name()]; found {
//				l.DebugCtx("Module already in cache", map[string]interface{}{
//					"context": "kbxutils",
//					"action":  "discoverModules",
//					"module":  info.Name(),
//					"path":    path,
//				})
//				return nil
//			}
//			// Add the module configuration to the cache
//			modulesConfig[info.Name()] = nil
//			return filepath.SkipDir
//		}
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	if len(modulesConfig) == 0 {
//		return nil, errors.New("no modules found")
//	}
//
//	return modulesConfig, nil
//}
