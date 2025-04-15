package settings

import (
	"github.com/faelmori/kubex-interfaces/settings"
	l "github.com/faelmori/logz"
)

// MetricsConfig is a struct that holds the configuration for metrics

type MetricsConfig struct {
	// Mutex for thread safety
	KubexThreading
	// Enabled enables metrics
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" gorm:"enabled,default:false"`
	// Address is the address for the metrics endpoint
	Address string `json:"address,omitempty" yaml:"address,omitempty" gorm:"address,default:localhost:9090"`
	// Path is the path for the metrics endpoint
	Path string `json:"path,omitempty" yaml:"path,omitempty" gorm:"path,default:/metrics"`
	// Interval is the interval for collecting metrics
	Interval int `json:"interval,omitempty" yaml:"interval,omitempty" gorm:"interval,default:60"`
}

// MetricsWatcherConfig is a struct that holds the configuration for metrics collection

type MetricsWatcherConfig struct {
	// Mutex for thread safety
	KubexThreading
	// Basic fields
	settings.KubexConfigBase
	// Metrics is the metrics configuration
	Metrics *MetricsConfig `json:"metrics,omitempty" yaml:"metrics,omitempty" gorm:"metrics"`
	// Upgrader is the upgrader for the metrics
	// Upgrader *MetricsUpgrader `json:"upgrader,omitempty" yaml:"upgrader,omitempty" gorm:"upgrader"`
	// LoggingModesConfig is the logging mode fields for the metrics
	LoggingModesConfig
}

// LoggingModesConfig is a struct that holds the base configuration for logging

type LoggingModesConfig struct {
	//// Mutex for thread safety
	//KubexThreading
	// Mode enables mode
	Mode bool `json:"auditMode,omitempty" yaml:"auditMode,omitempty" gorm:"auditMode,default:false"`
	// Level is the level (e.g., debug, info, warn, error)
	Level string `json:"auditLevel,omitempty" yaml:"auditLevel,omitempty" gorm:"auditLevel,default:info"`
	// Format is the format of the logs (e.g., json, text)
	Format string `json:"auditFormat,omitempty" yaml:"auditFormat,omitempty" gorm:"auditFormat,default:text"`
	// Output is the output destination for the logs (e.g., stdout, file)
	Output string `json:"auditOutput,omitempty" yaml:"auditOutput,omitempty" gorm:"auditOutput,default:stdout"`
	// Logger for logging
	Logger *l.Logger `json:"auditLogger,omitempty" yaml:"auditLogger,omitempty"`
	// Writers is a list of writers for the logs
	Writers []l.Writer `json:"auditWriters,omitempty" yaml:"auditWriters,omitempty" gorm:"auditWriters,default:stdout"`
	// FilePath is the path to the log file
	FilePath string `json:"auditFilePath,omitempty" yaml:"auditFilePath,omitempty" gorm:"auditFilePath,default:stdout"`
	// FileMaxSize is the maximum size of the log file in megabytes
	FileMaxSize int `json:"auditFileMaxSize,omitempty" yaml:"auditFileMaxSize,omitempty" gorm:"auditFileMaxSize,default:100"`
	// FileMaxBackups is the maximum number of backup files to keep
	FileMaxBackups int `json:"auditFileMaxBackups,omitempty" yaml:"auditFileMaxBackups,omitempty" gorm:"auditFileMaxBackups,default:10"`
	// FileMaxAge is the maximum age of the log files in days
	FileMaxAge int `json:"auditFileMaxAge,omitempty" yaml:"auditFileMaxAge,omitempty" gorm:"auditFileMaxAge,default:30"`
	// FileCompress enables compression for the log files
	FileCompress bool `json:"auditFileCompress,omitempty" yaml:"auditFileCompress,omitempty" gorm:"auditFileCompress,default:false"`
	// Metrics is the metrics configuration
	Metrics *MetricsConfig `json:"auditMetrics,omitempty" yaml:"auditMetrics,omitempty" gorm:"auditMetrics"`
}

// BaseLoggingConfig is a struct that holds the base configuration for logging

type BaseLoggingConfig struct {
	// Prefix is the prefix for the logs
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" gorm:"prefix,default:GoSpider"`
	// LogLevel is the logging level (e.g., debug, info, warn, error)
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty" gorm:"logLevel,default:info"`
	// LogFormat is the format of the logs (e.g., json, text)
	LogFormat string `json:"logFormat,omitempty" yaml:"logFormat,omitempty" gorm:"logFormat,default:text"`
	// Output is the output destination for the logs (e.g., stdout, file)
	Output string `json:"output,omitempty" yaml:"output,omitempty" gorm:"output,default:stdout"`
	// Logger for logging
	Logger *l.Logger `json:"logger,omitempty" yaml:"logger,omitempty"`
	// Writers is a list of writers for the logs
	Writers []l.Writer `json:"writers,omitempty" yaml:"writers,omitempty" gorm:"writers,default:stdout"`
}

// LoggingConfig is a struct that holds the configuration for logging

type LoggingConfig struct {
	KubexThreading

	// DebugMode enables debug mode
	DebugMode bool `json:"debugMode,omitempty" yaml:"debugMode,omitempty" gorm:"debugMode,default:false"`

	// BaseLogging is the base logging configuration
	BaseLogging BaseLoggingConfig

	// LoggingModes is the logging modes configuration
	LoggingModes LoggingModesConfig

	// Metrics is the metrics configuration
	Metrics *MetricsConfig `json:"metrics,omitempty" yaml:"metrics,omitempty" gorm:"metrics"`

	// Upgrader is the upgrader for the logging
	// Upgrader *LoggingUpgrader `json:"upgrader,omitempty" yaml:"upgrader,omitempty" gorm:"upgrader"`
}

// NewLoggingConfig creates a new LoggingConfig instance with default values
func NewLoggingConfig() *LoggingConfig {
	threading := NewThreading()
	return &LoggingConfig{
		KubexThreading: *threading,
		DebugMode:      false,
		BaseLogging: BaseLoggingConfig{
			Prefix:    "GoSpider",
			LogLevel:  "info",
			LogFormat: "text",
			Output:    "stdout",
			Logger:    nil,
			Writers:   []l.Writer{},
		},
		LoggingModes: LoggingModesConfig{
			settings.KubexThreading: NewThreading(),
			Metrics:                 settings.NewMetricsConfig(),
		},
		Metrics: NewMetricsConfig(),
	}
}
