package main

import (
	"github.com/ronaudinho/bigq/internal/handler"
	"github.com/ronaudinho/bigq/internal/service"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	defaultConf = &config.Config{
		Broker:          "amqp://guest:guest@localhost:5672/",
		DefaultQueue:    "machinery_tasks",
		ResultBackend:   "amqp://guest:guest@localhost:5672/",
		ResultsExpireIn: config.DefaultResultsExpireIn,
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "machinery_task",
			PrefetchCount: 3,
		},
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}
)

func main() {
	var conf *config.Config
	conf, err := config.NewFromYaml("config.yml", false)
	if err != nil {
		conf = defaultConf
	}
	svc := service.New(conf)
	hndlr := handler.New(svc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	setupRoutes(e, hndlr)
	e.Logger.Fatal(e.Start(":1323"))
}

func setupRoutes(e *echo.Echo, hndlr *handler.Handler) {
	e.POST("/argo", hndlr.RecvArgo)
	e.POST("/airflow", hndlr.RecvAirflow)
	e.GET("/argo", hndlr.SendArgo)
	e.GET("/airflow", hndlr.SendAirflow)
}
