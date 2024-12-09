//go:build !prod

package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kingwrcy/moments/db"
	_ "github.com/kingwrcy/moments/docs"
	"github.com/kingwrcy/moments/handler"
	"github.com/kingwrcy/moments/log"
	myMiddleware "github.com/kingwrcy/moments/middleware"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	_ "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var version string
var commitId string

func newEchoEngine(_ do.Injector) (*echo.Echo, error) {
	e := echo.New()
	return e, nil
}

// @title	Moments API
// @version	0.2.1
func main() {
	injector := do.New()
	var cfg vo.AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Printf("读取配置文件异常:%s", err)
		return
	}

	if version == "" {
		version = "unknown"
	}

	if commitId == "" {
		commitId = "unknown"
	}

	do.ProvideValue(injector, &cfg)
	do.Provide(injector, log.NewLogger)

	myLogger := do.MustInvoke[zerolog.Logger](injector)

	myLogger.Info().Msgf("version = %s", version)
	myLogger.Info().Msgf("commitId = %s", commitId)

	handleEmptyConfig(myLogger, &cfg)
	cfg.Version = version
	cfg.CommitId = commitId

	do.Provide(injector, db.NewDB)
	do.Provide(injector, newEchoEngine)
	do.Provide(injector, handler.NewBaseHandler)

	tx := do.MustInvoke[*gorm.DB](injector)

	e := do.MustInvoke[*echo.Echo](injector)
	e.Use(myMiddleware.Auth(injector))

	setupRouter(injector)

	migrateTo3(tx, myLogger)
	e.HideBanner = true
	myLogger.Info().Msgf("服务端启动成功,监听:%d端口...", cfg.Port)
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		myLogger.Fatal().Msgf("服务启动失败,错误原因:%s", err)
	}
}
