package main

import (
	"embed"
	"log"

	"git.gnous.eu/gnouseu/plakken/internal/config"
	"git.gnous.eu/gnouseu/plakken/internal/database"
	"git.gnous.eu/gnouseu/plakken/internal/httpServer"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed static
	static embed.FS
)

func main() {
	initConfig := config.GetConfig()
	dbConfig := database.InitDB(initConfig.RedisAddress, initConfig.RedisUser, initConfig.RedisPassword, initConfig.RedisDB)
	db := database.ConnectDB(dbConfig)
	err := database.Ping(db)
	if err != nil {
		log.Fatal(err)
	}

	serverConfig := httpServer.ServerConfig{
		HTTPServer: httpServer.Config(initConfig.ListenAddress),
		UrlLength:  initConfig.UrlLength,
		DB:         db,
		Static:     static,
		Templates:  templates,
	}

	serverConfig.Server()
}
