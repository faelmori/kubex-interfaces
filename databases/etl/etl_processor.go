package etl

import (
	"fmt"
	"github.com/faelmori/kubex-interfaces/databases/drivers"
	"github.com/faelmori/kubex-interfaces/databases/types"
	tl "github.com/faelmori/kubex-interfaces/tools"
	t "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"reflect"
	"time"
)

type DBProcessor struct {
	InputChannel  t.IChannel[types.Config, int]
	OutputChannel t.IChannel[TableHandler, int]
	Logger        l.Logger
}

func NewETLProcessor() *DBProcessor {
	return &DBProcessor{
		InputChannel:  tl.NewChannel[types.Config, int]("InputChannel", nil, 10),
		OutputChannel: tl.NewChannel[TableHandler, int]("OutputChannel", nil, 10),
		Logger:        l.GetLogger("DBProcessor"),
	}
}

func (e *DBProcessor) StartProcessing() {
	go func() {
		chanCfg, tpyCfg, errCfg := e.InputChannel.Monitor()
		if errCfg != nil {
			e.Logger.ErrorCtx("Failed to process ETL: "+errCfg.Error(), nil)
			return
		}
		for cfgObj := range chanCfg {
			if reflect.TypeOf(cfgObj) != tpyCfg {
				e.Logger.ErrorCtx("Failed to process ETL: Type mismatch", nil)
				return
			}
			config := cfgObj.(types.Config)
			db, err := drivers.ConnectDB(config)
			if err != nil {
				e.Logger.ErrorCtx("Failed to connect to database: "+err.Error(), nil)
				continue
			}
			if table, err := ExtractDataWithTableHandler(db, config); err != nil {
				e.Logger.ErrorCtx("Failed to process ETL: "+err.Error(), nil)
				continue
			} else {
				if sendErr := e.OutputChannel.Send(*table); sendErr != nil {
					return
				}
			}
		}
	}()
}

func (e *DBProcessor) MonitorProgress(interval time.Duration) {
	go func() {
		// Monitorar o progresso do canal de saída
		lst, tpy := e.OutputChannel.GetChan()
		for {
			select {
			case data := <-lst:
				if reflect.TypeOf(data) != tpy {
					e.Logger.ErrorCtx("Failed to process ETL: Type mismatch", nil)
					continue
				} else {
					table := data.(TableHandler)
					e.Logger.InfoCtx(fmt.Sprintf("Processed %d rows", len(table.GetRows())), nil)
				}
			}
		}
	}()
}
