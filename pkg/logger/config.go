package logger

type Type string

const (
	Std     Type = "std"
	Fluent  Type = "fluent"
	Zap     Type = "zap"
	Logrus  Type = "logrus"
	Aliyun  Type = "aliyun"
	Tencent Type = "tencent"
)

type Config struct {
	LoggerType  Type
	Filename    string
	Level       string
	MaxSize     int32
	MaxAge      int32
	MaxBackups  int32
	Stdout      bool
	SimpleTrace bool
}

type Option func(*Config)

func WithLoggerType(t Type) Option {
	return func(c *Config) {
		c.LoggerType = t
	}
}

func WithFile(Filename string) Option {
	return func(c *Config) {
		c.Filename = Filename
	}
}
func WithLevel(level string) Option {
	return func(c *Config) {
		c.Level = level
	}
}
func WithMaxSize(maxSize int32) Option {
	return func(c *Config) {
		c.MaxSize = maxSize
	}
}
func WithMaxAge(maxAge int32) Option {
	return func(c *Config) {
		c.MaxAge = maxAge
	}
}
func WithMaxBackups(maxBackups int32) Option {
	return func(c *Config) {
		c.MaxBackups = maxBackups
	}
}
func WithStdout(stdout bool) Option {
	return func(c *Config) {
		c.Stdout = stdout
	}
}
func WithSimpleTrace(simpleTrace bool) Option {
	return func(c *Config) {
		c.SimpleTrace = simpleTrace
	}
}

var defaultConfig = &Config{
	LoggerType:  Std,
	Filename:    "info.log",
	Level:       "info",
	MaxSize:     1024,
	MaxAge:      7,
	MaxBackups:  3,
	Stdout:      true,
	SimpleTrace: false,
}
