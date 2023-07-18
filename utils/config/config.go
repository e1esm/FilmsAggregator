package config

import (
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	confFile = "conf.yml"
)

type Config struct {
	Aggregator struct {
		Name    string `yaml:"name"`
		Port    int    `yaml:"port"`
		Address string `yaml:"address"`
	} `yaml:"aggregator"`
	Reindexer struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	} `yaml:"reindexer"`
	CacheTime string `yaml:"cache_time"`
	Postgres  struct {
		ContainerName string `yaml:"container_name"`
		Port          int    `yaml:"port"`
		User          string `yaml:"user"`
		DatabaseName  string `yaml:"database_name"`
		Password      string `yaml:"password"`
		Connections   int    `yaml:"connections"`
	} `yaml:"postgres"`
}

func NewConfig() *Config {
	fileContent, err := os.ReadFile(confFile)
	if err != nil {
		logger.Logger.Error("Couldn't have read the file",
			zap.String("err", err.Error()),
			zap.String("filename", confFile))
		return nil
	}
	cfg := &Config{}
	err = yaml.Unmarshal(fileContent, cfg)
	if err != nil {
		logger.Logger.Error("Couldn't have unmarshalled content of file",
			zap.String("filename", confFile), zap.String("err", err.Error()))
		return nil
	}
	return cfg

}
