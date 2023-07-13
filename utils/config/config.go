package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Aggregator struct {
		Name string `yaml:"name"`
	} `yaml:"aggregator"`
	DBCOnnection struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	} `yaml:"reindexer"`
	CacheTime string `yaml:"cache_time"`
}

func NewConfig() *Config {
	fileContent, err := os.ReadFile("conf.yml")
	if err != nil {
		log.Println("Couldn't have read file")
		return nil
	}
	cfg := &Config{}
	err = yaml.Unmarshal(fileContent, cfg)
	if err != nil {
		log.Println("Couldn't have unmarshalled utils")
		return nil
	}
	return cfg

}
