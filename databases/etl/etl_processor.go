package etl

import (
	"encoding/json"
	"fmt"
	"github.com/faelmori/kubex-interfaces/databases/drivers"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	tl "github.com/faelmori/kubex-interfaces/tools"
	tp "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"os"
	"reflect"
	"sync"
	"time"
)

type ETLProcessor struct {
	mu            sync.Mutex
	wg            sync.WaitGroup
	Logger        l.Logger
	InputChannel  tp.IChannel[t.Config, int]
	OutputChannel tp.IChannel[TableHandler, int]
}

func NewETLProcessor(logger l.Logger) *ETLProcessor {
	if logger == nil {
		logger = l.NewLogger("Kubex ETL Processor")
	}
	return &ETLProcessor{
		InputChannel:  tl.NewChannel[t.Config, int]("InputChannel", nil, 10),
		OutputChannel: tl.NewChannel[TableHandler, int]("OutputChannel", nil, 10),
		Logger:        logger,
	}
}

func (e *ETLProcessor) StartProcessing() {
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
			//config := cfgObj.(t.Config)
			config := cfgObj.(t.DBConfigBase)
			_, err := drivers.ConnectDB(&config)
			if err != nil {
				e.Logger.ErrorCtx("Failed to connect to database: "+err.Error(), nil)
				continue
			}
			if table, err := ExtractDataWithTableHandler(nil, t.Config{}); err != nil {
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

func (e *ETLProcessor) MonitorProgress(interval time.Duration) {
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

func (e *ETLProcessor) LoadJobFromFile(filePath string) (*t.VJob, error) {
	fileData, fileDataErr := os.ReadFile(filePath)
	if fileDataErr != nil {
		l.ErrorCtx("failed to load file: "+fileDataErr.Error(), map[string]interface{}{})
		return nil, fileDataErr
	}

	var job t.VJob
	if unmarshalErr := json.Unmarshal(fileData, &job); unmarshalErr != nil {
		l.ErrorCtx("failed to unmarshal file data: "+unmarshalErr.Error(), map[string]interface{}{})
		return nil, unmarshalErr
	}

	return &job, nil
}
