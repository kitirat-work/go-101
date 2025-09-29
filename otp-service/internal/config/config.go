package config

import (
	"os"
	"reflect"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`
	MySqlConfig
}

type MySqlConfig struct {
	Url             string `env:"MYSQL_CLIENT_URL" envDefault:"mysql://myuser:mypassword@tcp(127.0.0.1:6033)/mydatabase"`
	ConnMaxLifetime string `env:"MYSQL_CLIENT_CONN_MAX_LIFETIME" envDefault:"3600"`
	MaxOpenConns    string `env:"MYSQL_CLIENT_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    string `env:"MYSQL_CLIENT_MAX_IDLE_CONNS" envDefault:"100"`
}

func NewConfig() *Config {
	cfg := &Config{}
	reflectStruct(cfg)
	return cfg
}

func reflectStruct(s any) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		envTag := field.Tag.Get("env")
		envDefaultTag := field.Tag.Get("envDefault")

		if val.Field(i).Kind() == reflect.Struct {
			reflectStruct(val.Field(i).Addr().Interface())
			continue
		}

		if envTag != "" {
			envValue, ok := os.LookupEnv(envTag)
			if ok {
				val.Field(i).SetString(envValue)
			} else {
				val.Field(i).SetString(envDefaultTag)
			}
		}
	}
}
