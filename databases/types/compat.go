package types

import (
	"fmt"
	"github.com/faelmori/logz"
)

const batchSize = 1000

type Field map[string]string
type Fields map[string][]Field

type Config struct {
	SourceType                  string           `json:"sourceType"`
	SourceConnectionString      string           `json:"sourceConnectionString"`
	SourceTable                 string           `json:"sourceTable"`
	DestinationType             string           `json:"destinationType"`
	DestinationConnectionString string           `json:"destinationConnectionString"`
	DestinationTable            string           `json:"destinationTable"`
	SQLQuery                    string           `json:"sqlQuery"`
	OutputPath                  string           `json:"outputPath"`
	OutputFormat                string           `json:"outputFormat"`
	Transformations             []Transformation `json:"transformations"`
	NeedCheck                   bool             `json:"needCheck"`
	CheckMethod                 string           `json:"checkMethod"`
	Joins                       []Join           `json:"joins"`
	Where                       string           `json:"where"`
	OrderBy                     string           `json:"orderBy"`
	Triggers                    []Trigger        `json:"triggers"`
	LogTable                    string           `json:"logTable"`
	SyncInterval                string           `json:"syncInterval"`
	KafkaURL                    string           `json:"kafkaURL"`
	KafkaTopic                  string           `json:"kafkaTopic"`
	KafkaGroupID                string           `json:"kafkaGroupID"`
	PrimaryKey                  string           `json:"primaryKey"`
	UpdateKey                   string           `json:"updateKey"`
}
type Transformation struct {
	SourceField      string `json:"sourceField"`
	DestinationField string `json:"destinationField"`
	Operation        string `json:"operation"`
	SPath            string `json:"sPath"`
	DPath            string `json:"dPath"`
	Type             string `json:"type"`
}
type Join struct {
	Table     string `json:"table"`
	Condition string `json:"condition"`
	JoinType  string `json:"joinType"`
}
type Trigger struct {
	Name      string `json:"name"`
	Table     string `json:"table"`
	Event     string `json:"event"`
	Statement string `json:"statement"`
}
type VendorSqlTypeMap struct {
	SourceType string
	TargetType string
	Fallback   string
}
type VendorSqlTypeMapList []VendorSqlTypeMap
type VendorSqlMapping struct {
	driver  string
	mapping VendorSqlTypeMapList
}
type VendorSqlMappingList []VendorSqlMapping

var vAendorMappingList = VendorSqlMappingList{
	{
		driver: "sqlite3",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "REAL", "REAL"},
			{"VARCHAR", "TEXT", "TEXT"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "INTEGER", "INTEGER"},
			{"DECIMAL", "REAL", "REAL"},
			{"VARCHAR2", "TEXT", "TEXT"},
			{"DATE", "TEXT", "TEXT"},
			{"DATETIME", "TEXT", "TEXT"},
			{"TIMESTAMP", "TEXT", "TEXT"},
			{"BOOLEAN", "INTEGER", "INTEGER"},
			{"BLOB", "BLOB", "BLOB"},
			{"CLOB", "CLOB", "CLOB"},
		},
	},
	{
		driver: "sqlite",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "REAL", "REAL"},
			{"VARCHAR", "TEXT", "TEXT"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "REAL", "REAL"},
			{"DECIMAL", "REAL", "REAL"},
			{"VARCHAR2", "TEXT", "TEXT"},
			{"DATE", "TEXT", "TEXT"},
			{"DATETIME", "TEXT", "TEXT"},
			{"TIMESTAMP", "TEXT", "TEXT"},
			{"BOOLEAN", "INTEGER", "INTEGER"},
			{"BLOB", "BLOB", "BLOB"},
			{"CLOB", "CLOB", "CLOB"},
			{"REAL", "REAL", "REAL"},
			{"FLOAT", "REAL", "REAL"},
		},
	},
	{
		driver: "postgres",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "NUMERIC", "NUMERIC"},
			{"VARCHAR", "VARCHAR", "VARCHAR"},
			{"NVARCHAR", "VARCHAR", "VARCHAR"},
			{"CHAR", "TEXT", "TEXT"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "INTEGER", "INTEGER"},
			{"SMALLINT", "INTEGER", "INTEGER"},
			{"DECIMAL", "NUMERIC", "NUMERIC"},
			{"VARCHAR2", "VARCHAR", "VARCHAR"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "TIMESTAMP", "TIMESTAMP"},
			{"TIMESTAMP", "TIMESTAMP", "TIMESTAMP"},
			{"BOOLEAN", "BOOLEAN", "BOOLEAN"},
			{"BLOB", "BYTEA", "BYTEA"},
			{"CLOB", "TEXT", "TEXT"},
			{"REAL", "REAL", "REAL"},
			{"FLOAT", "REAL", "REAL"},
		},
	},
	{
		driver: "mysql",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "INT", "INT"},
			{"VARCHAR", "VARCHAR", "VARCHAR"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "INT", "INT"},
			{"DECIMAL", "DECIMAL", "DECIMAL"},
			{"VARCHAR2", "VARCHAR", "VARCHAR"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "DATETIME", "DATETIME"},
			{"TIMESTAMP", "TIMESTAMP", "TIMESTAMP"},
			{"BOOLEAN", "TINYINT", "TINYINT"},
			{"BLOB", "BLOB", "BLOB"},
			{"CLOB", "TEXT", "TEXT"},
			{"REAL", "REAL", "REAL"},
			{"FLOAT", "FLOAT", "FLOAT"},
		},
	},
	{
		driver: "oracle",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "NUMBER", "NUMBER"},
			{"VARCHAR", "VARCHAR2", "VARCHAR2"},
			{"TEXT", "CLOB", "CLOB"},
			{"INT", "NUMBER", "NUMBER"},
			{"DECIMAL", "NUMBER", "NUMBER"},
			{"VARCHAR2", "VARCHAR2", "VARCHAR2"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "TIMESTAMP", "TIMESTAMP"},
			{"TIMESTAMP", "TIMESTAMP", "TIMESTAMP"},
			{"BOOLEAN", "NUMBER", "NUMBER"},
			{"BLOB", "BLOB", "BLOB"},
			{"CLOB", "CLOB", "CLOB"},
			{"REAL", "REAL", "REAL"},
		},
	},
	{
		driver: "sqlserver",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "DECIMAL", "DECIMAL"},
			{"VARCHAR", "VARCHAR", "VARCHAR"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "INT", "INT"},
			{"DECIMAL", "DECIMAL", "DECIMAL"},
			{"VARCHAR2", "VARCHAR", "VARCHAR"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "DATETIME", "DATETIME"},
			{"TIMESTAMP", "DATETIME", "DATETIME"},
			{"BOOLEAN", "BIT", "BIT"},
			{"BLOB", "VARBINARY", "VARBINARY"},
			{"CLOB", "TEXT", "TEXT"},
			{"REAL", "REAL", "REAL"},
			{"FLOAT", "FLOAT", "FLOAT"},
		},
	},
	{
		driver: "mssql",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "DECIMAL", "DECIMAL"},
			{"VARCHAR", "VARCHAR", "VARCHAR"},
			{"TEXT", "TEXT", "TEXT"},
			{"INT", "INT", "INT"},
			{"DECIMAL", "DECIMAL", "DECIMAL"},
			{"VARCHAR2", "VARCHAR", "VARCHAR"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "DATETIME", "DATETIME"},
			{"TIMESTAMP", "DATETIME", "DATETIME"},
			{"BOOLEAN", "BIT", "BIT"},
			{"BLOB", "VARBINARY", "VARBINARY"},
			{"CLOB", "TEXT", "TEXT"},
			{"REAL", "REAL", "REAL"},
			{"FLOAT", "FLOAT", "FLOAT"},
		},
	},
	{
		driver: "godror",
		mapping: VendorSqlTypeMapList{
			{"NUMBER", "NUMBER", "NUMBER"},
			{"VARCHAR", "VARCHAR2", "VARCHAR2"},
			{"TEXT", "CLOB", "CLOB"},
			{"INT", "NUMBER", "NUMBER"},
			{"DECIMAL", "NUMBER", "NUMBER"},
			{"VARCHAR2", "VARCHAR2", "VARCHAR2"},
			{"DATE", "DATE", "DATE"},
			{"DATETIME", "TIMESTAMP", "TIMESTAMP"},
			{"TIMESTAMP", "TIMESTAMP", "TIMESTAMP"},
			{"BOOLEAN", "NUMBER", "NUMBER"},
			{"BLOB", "BLOB", "BLOB"},
			{"CLOB", "CLOB", "CLOB"},
			{"REAL", "REAL", "REAL"},
		},
	},
}

func GetVendorSqlTypeMap(driver string) VendorSqlTypeMapList {
	for _, mapping := range vAendorMappingList {
		if mapping.driver == driver {
			return mapping.mapping
		}
	}
	logz.ErrorCtx(fmt.Sprintf("No mapping found for driver %s", driver), map[string]interface{}{})
	return nil
}
func GetVendorSqlType(driver, sourceType string) string {
	mapping := GetVendorSqlTypeMap(driver)
	if mapping == nil {
		logz.ErrorCtx(fmt.Sprintf("No mapping found for driver %s", driver), map[string]interface{}{})
		return ""
	}
	for _, mapItem := range mapping {
		if mapItem.SourceType == sourceType {
			return mapItem.TargetType
		}
	}
	logz.ErrorCtx(fmt.Sprintf("No mapping found for source type %s", sourceType), map[string]interface{}{})
	return ""
}
