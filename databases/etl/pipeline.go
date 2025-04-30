package etl

import l "github.com/faelmori/logz"

type ETLPipeline struct {
	ConfigManager   ConfigManager
	DataExtractor   DataExtractor
	DataTransformer DataTransformer
	DataLoader      DataLoader
	Logger          l.Logger
}

func (etl *ETLPipeline) RunPipeline() error {
	handler, err := etl.DataExtractor.ExtractData(etl.ConfigManager.Config.SQLQuery)
	if err != nil {
		return err
	}

	etl.DataTransformer.ApplyTransformations(handler)

	if err := etl.DataLoader.LoadData(handler); err != nil {
		return err
	}

	etl.Logger.InfoCtx("ETL process completed successfully!", nil)
	return nil
}
