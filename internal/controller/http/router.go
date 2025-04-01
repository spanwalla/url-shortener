package http

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/spanwalla/url-shortener/internal/controller/http/api/get_alias"
	"github.com/spanwalla/url-shortener/internal/controller/http/api/post_root"
	"github.com/spanwalla/url-shortener/internal/service"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.CORS())
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	getAlias := get_alias.New(services.Expander)
	postRoot := post_root.New(services.Shortener)

	handler.GET("/swagger/*", echoSwagger.WrapHandler)
	handler.GET("/:alias", getAlias.Handle)
	handler.POST("/", postRoot.Handle)
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("http - setLogsFile - os.OpenFile: %v", err)
	}

	return file
}
