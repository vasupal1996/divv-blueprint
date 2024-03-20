package config

import (
	"fmt"
)

type Config struct {
	AppConfig       *AppConfig       `mapstructure:"app_config"`
	MongoDBConfig   *MongoDBConfig   `mapstructure:"mongo_db_config"`
	WebServerConfig *WebServerConfig `mapstructure:"web_server_config"`
	RouterConfig    *RouterConfig    `mapstructure:"router_config"`
	SentryConfig    *SentryConfig    `mapstructure:"sentry_config"`
}

type AppConfig struct {
	ServiceConfig *ServiceConfig `mapstructure:"service_config"`
}

type ServiceConfig struct {
	DemoServiceConfig *DemoServiceConfig `mapstructure:"demo_service_config"`
}

type RouterConfig struct {
	EnableSentry bool `mapstructure:"enable_sentry"`
}

type DemoServiceConfig struct {
	SomeAdditionalData string `mapstructure:"some_additional_data"`
}

/*
DATABASE RELATED CONFIG
*/

type MongoDBConfig struct {
	Scheme     string `mapstructure:"scheme"`
	Host       string `mapstructure:"host"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	ReplicaSet string `mapstructure:"replica_set"`
	AppName    string `mapstructure:"app_name"`
	ReadPref   string `mapstructure:"read_pref"`
}

func (d *MongoDBConfig) ConnectionURL() string {
	url := fmt.Sprintf("%s://", d.Scheme)
	if d.Username != "" && d.Password != "" {
		url += fmt.Sprintf("%s:%s@", d.Username, d.Password)
	}
	url += fmt.Sprintf("%s/?", d.Host)
	if d.ReplicaSet != "" {
		url += fmt.Sprintf("replicaSet=%s", d.ReplicaSet)
	}
	if d.AppName != "" {
		url += fmt.Sprintf("appName=%s", d.AppName)
	}
	return url
}

/*
WEB SERVER RELATED CONFIG
*/

type WebServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

/*
SENTRY AND TRACKING CONFIG
*/

type SentryConfig struct {
	EnableSentry bool   `mapstructure:"enable_sentry"`
	Host         string `mapstructure:"dsn"`
}
