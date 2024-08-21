package main

import "project-manager-go/data"

type ConfigServer struct {
	URL  string
	Port string
	Cors []string
}

type AppConfig struct {
	Server     ConfigServer
	DB         data.DBConfig
	Demodata   string
	Binarydata string
}
