package main

import (
	"fmt"
	"wallet-service/internal/command"
	"wallet-service/internal/config"
	"wallet-service/internal/utils"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	executor := command.NewCommandExecutor(viperConfig, db)
	jwt := utils.NewJWT(viperConfig)
	validate := config.NewValidator(viperConfig)
	router := config.NewGin(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Router:   router,
		Log:      log,
		JWT:      jwt,
		Validate: validate,
		Config:   viperConfig,
	})

	if !executor.Execute(log) {
		return
	}

	webPort := viperConfig.GetInt("PORT")
	err := router.Run(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
