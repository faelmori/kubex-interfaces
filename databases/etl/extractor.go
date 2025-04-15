package etl

import (
	"database/sql"
	"fmt"
	"github.com/elgris/sqrl"
	d "github.com/faelmori/kubex-interfaces/databases/drivers"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	l "github.com/faelmori/logz"
	"strings"
)

type DataExtractor struct {
	Config ConfigManager
	Logger l.Logger
}

func (de *DataExtractor) ExtractData(query string) (*TableHandler, error) {
	db, err := d.ConnectDB(de.Config.Config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	data, err := de.MapRowsToData(rows, columns)
	if err != nil {
		return nil, err
	}

	return &TableHandler{Columns: columns, Data: t.ConvertDataToRows(data)}, nil
}

func (de *DataExtractor) MapRowsToData(rows *sql.Rows, columns []string) ([]t.Data, error) {
	var data []t.Data
	for rows.Next() {
		rowData := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range rowData {
			rowPointers[i] = &rowData[i]
		}

		if scanErr := rows.Scan(rowPointers...); scanErr != nil {
			return nil, scanErr
		}

		row := t.Data{}
		for i, colName := range columns {
			row[colName] = rowData[i]
		}
		data = append(data, row)
	}

	return data, nil
}

func ExtractDataWithTypes(dbSQL *sql.DB, config t.Config) ([]t.Data, map[string]string, error) {
	var db *sql.DB
	var dbErr error
	if dbSQL == nil {
		db, dbErr = sql.Open(config.SourceType, config.SourceConnectionString)
		if dbErr != nil {
			l.ErrorCtx("Failed to connect to source database: "+dbErr.Error(), map[string]interface{}{})
			return nil, nil, dbErr
		}
	} else {
		db = dbSQL
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	l.InfoCtx("Starting data extraction", map[string]interface{}{})

	var rows *sql.Rows
	var SQLQueryArgs []interface{}
	var rowsErr error
	var buildQueryErr error

	if config.SQLQuery == "" {
		var fields []string
		var transformationsList []t.Transformation
		transformationsList = config.Transformations
		for i, trans := range transformationsList {
			fields = append(fields, trans.SourceField)
			if trans.Type == "" {
				transformationsList[i].Type = "string"
			}
		}
		config.SQLQuery, SQLQueryArgs, buildQueryErr = BuilExtractdQuery(config, fields)
		if buildQueryErr != nil {
			l.ErrorCtx("Failed to build query: "+buildQueryErr.Error(), map[string]interface{}{})
			return nil, nil, buildQueryErr
		}
	}

	//logz.DebugLog("Running query: "+config.SQLQuery, map[string]interface{}{})

	if len(SQLQueryArgs) > 0 {
		rows, rowsErr = db.Query(config.SQLQuery, SQLQueryArgs...)
	} else {
		rows, rowsErr = db.Query(config.SQLQuery)
	}

	if rowsErr != nil {
		l.ErrorCtx("Failed on query execution: "+rowsErr.Error(), map[string]interface{}{})
		return nil, nil, rowsErr
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var data []t.Data
	columns, columnsErr := rows.Columns()
	if columnsErr != nil {
		l.ErrorCtx("Failed to get columns: "+columnsErr.Error(), map[string]interface{}{})
		return nil, nil, columnsErr
	}

	columnTypes, columnTypesErr := rows.ColumnTypes()
	if columnTypesErr != nil {
		l.ErrorCtx("Failed trying to get column types: "+columnTypesErr.Error(), map[string]interface{}{})
		return nil, nil, columnTypesErr
	}

	columnTypeMap := make(map[string]string)
	for i, colType := range columnTypes {
		columnTypeMap[columns[i]] = colType.DatabaseTypeName()
	}

	for rows.Next() {
		rowData := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range rowData {
			rowPointers[i] = &rowData[i]
		}

		if scanErr := rows.Scan(rowPointers...); scanErr != nil {
			l.ErrorCtx("Failed to scan row data: "+scanErr.Error(), map[string]interface{}{})
			return nil, nil, scanErr
		}

		row := t.Data{}
		for i, colName := range columns {
			row[colName] = rowData[i]
		}
		data = append(data, row)
	}

	return data, columnTypeMap, nil
}

func ExtractDataWithTableHandler(dbSQL *sql.DB, config t.Config) (*TableHandler, error) {
	data, columnTypeMap, err := ExtractDataWithTypes(dbSQL, config)
	if err != nil {
		return nil, err
	}

	// Criar o handler com base nos dados extraídos
	handler := &TableHandler{
		Types:   t.ConvertMapToTypeArray(columnTypeMap),
		Headers: t.ExtractHeaders(columnTypeMap),
		Rows:    t.ConvertDataToRows(data),
	}

	return handler, nil
}

func ExtractData(dbSQL *sql.DB, config t.Config) ([]t.Data, []string, error) {
	if config.SQLQuery == "" {
		l.ErrorCtx("query SQL não informada", map[string]interface{}{})
		return nil, nil, fmt.Errorf("query SQL não informada")
	}

	var db *sql.DB
	var dbErr error
	if dbSQL == nil {
		db, dbErr = sql.Open(config.SourceType, config.SourceConnectionString)
		if dbErr != nil {
			l.ErrorCtx(fmt.Sprintf("falha ao conectar ao banco de dados: %v", dbErr), map[string]interface{}{})
			return nil, nil, dbErr
		}
	} else {
		db = dbSQL
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, queryErr := db.Query(config.SQLQuery)
	if queryErr != nil {
		l.ErrorCtx(fmt.Sprintf("falha ao executar a query SQL: %v", queryErr), map[string]interface{}{})
		return nil, nil, queryErr
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var data []t.Data
	columns, columnsErr := rows.Columns()
	if columnsErr != nil {
		l.ErrorCtx(fmt.Sprintf("falha ao obter colunas: %v", columnsErr), map[string]interface{}{})
		return nil, nil, columnsErr
	}

	for rows.Next() {
		rowData := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range rowData {
			rowPointers[i] = &rowData[i]
		}

		if scanErr := rows.Scan(rowPointers...); scanErr != nil {
			l.ErrorCtx(fmt.Sprintf("falha ao escanear os dados da linha: %v", scanErr), map[string]interface{}{})
			return nil, nil, scanErr
		}

		row := t.Data{}
		for i, colName := range columns {
			row[colName] = rowData[i]
		}
		data = append(data, row)
	}

	if config.OutputPath != "" {
		saveDataErr := SaveData(config.OutputPath, data, config.OutputFormat)
		if saveDataErr != nil {
			l.ErrorCtx("Failed to save data: "+saveDataErr.Error(), map[string]interface{}{})
		}
	}

	return data, columns, nil
}

func BuilExtractdQuery(config t.Config, fields []string) (string, []interface{}, error) {
	query := sqrl.Select(fields...).From(config.SourceTable)

	for _, join := range config.Joins {
		switch strings.ToUpper(join.JoinType) {
		case "INNER":
			query = query.Join(join.Table + " ON " + join.Condition)
		case "LEFT":
			query = query.LeftJoin(join.Table + " ON " + join.Condition)
		case "RIGHT":
			query = query.RightJoin(join.Table + " ON " + join.Condition)
		default:
			return "", nil, fmt.Errorf("tipo de join desconhecido: %s", join.JoinType)
		}
	}

	if config.Where != "" {
		query = query.Where(config.Where)
	}

	if config.OrderBy != "" {
		query = query.OrderBy(config.OrderBy)
	}

	return query.ToSql()
}

func GetDataTableHandlerFromQuery(sourceType, sourceConnectionString, sqlQuery string) (*TableHandler, error) {
	db, err := sql.Open(sourceType, sourceConnectionString)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados de origem: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a consulta SQL: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter colunas: %w", err)
	}

	var data [][]string
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha: %w", err)
		}

		var row []string
		for _, value := range values {
			row = append(row, fmt.Sprintf("%v", value))
		}
		data = append(data, row)
	}

	return &TableHandler{Columns: columns, Data: data}, nil
}
