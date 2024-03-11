package config

import (
	viper "github.com/spf13/viper"
)

type App struct {
	Port int `json:"port"`
}

type Database struct {
	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Jaeger struct {
	ServiceName string `json:"service_name,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
}

type JWT struct {
	SecretKey string `json:"secret_key,omitempty"`
}

type ConfigApp struct {
	App      *App      `json:"app"`
	Database *Database `json:"database"`
	Jaeger   *Jaeger   `json:"jaeger"`
	Jwt      *JWT      `json:"jwt"`
}

func NewConfigApp() IConfig {
	viper := viper.New()
	viper.SetConfigFile("config.json")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	cfg := ConfigApp{
		App: &App{Port: viper.GetInt("app.port")},
		Database: &Database{
			Host:     viper.GetString("database.host"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Port:     viper.GetInt("database.port"),
			Name:     viper.GetString("database.name"),
		},
		Jaeger: &Jaeger{
			ServiceName: viper.GetString("jaeger.service_name"),
			Host:        viper.GetString("jaeger.host"),
			Port:        viper.GetInt("jaeger.port"),
		},
		Jwt: &JWT{SecretKey: viper.GetString("jwt.secret_key")},
	}

	return &cfg
}

func (c *ConfigApp) Config() *ConfigApp {
	return c
}
