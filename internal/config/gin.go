package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewGin(config *viper.Viper) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.SetTrustedProxies(nil)
	return engine
}
