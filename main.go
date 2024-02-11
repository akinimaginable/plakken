package main

import (
	"git.gnous.eu/gnouseu/plakken/internal/config"
	"git.gnous.eu/gnouseu/plakken/internal/database"
	"git.gnous.eu/gnouseu/plakken/internal/httpServer"
)

func main() {
	initConfig := config.GetConfig()
	dbConfig := database.InitDB(initConfig.RedisAddress, initConfig.RedisUser, initConfig.RedisPassword, initConfig.RedisDB)
	db := database.ConnectDB(dbConfig)

	serverConfig := httpServer.ServerConfig{
		HTTPServer: httpServer.Config(initConfig.ListenAddress),
		UrlLength:  initConfig.UrlLength,
		DB:         db,
	}

	serverConfig.Server()
}
