package config

import viper "github.com/spf13/viper"

type App struct {
	Port int `json:"port"`
}

type ConfigApp struct {
	App *App
}

func NewConfigApp() IConfig {
	viper := viper.New()
	viper.SetConfigFile("config.json")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	cfg := ConfigApp{
		App: &App{Port: viper.GetInt("app.port")},
	}

	return &cfg
}

func (c *ConfigApp) Config() *ConfigApp {
	return c
}
