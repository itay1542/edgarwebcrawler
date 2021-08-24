package main

import "github.com/itay1542/edgarwebcrawler/DAL"

type Config struct {
	AlphaVantage struct {
		Host   string `yaml:"host"`
		ApiKey string `yaml:"apiKey"`
	} `yaml:"alphaVantage"`
	Filter struct {
		StockExchanges []DAL.StockExchange `yaml:"stockExchanges"`
	} `yaml:"filter"`
	DB struct {
		Host     string `yaml:"host"`
		Port     uint   `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}
