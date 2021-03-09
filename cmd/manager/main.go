package main

import (
	"github.com/ronaudinho/bigq/pkg/task"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
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
	defaultConsumerTag = "machinery_worker"
)

func main() {
	err := worker()
	if err != nil {
		log.FATAL.Fatal(err)
	}
}

func startServer() (*machinery.Server, error) {
	var conf *config.Config
	conf, err := config.NewFromYaml("config.yml", false)
	if err != nil {
		conf = defaultConf
	}

	server, err := machinery.NewServer(conf)
	if err != nil {
		return nil, err
	}

	return server, server.RegisterTasks(processTasks())
}

// worker spawns unlimited concurrent worker with defaultConsumerTag
// TODO should probably refactor this to set up different exchanges
// and accept unique consumerTag
func worker() error {
	errorsChan := make(chan error)
	for k := range processTasks() {
		server, err := startServer()
		if err != nil {
			return err
		}

		worker := server.NewCustomQueueWorker(k+"_consumer", 0, k)
		errorhandler := func(err error) {
			log.ERROR.Println("error:", err)
		}

		pretaskhandler := func(sig *tasks.Signature) {
			log.INFO.Println("starting:", sig.Name)
		}

		posttaskhandler := func(sig *tasks.Signature) {
			log.INFO.Printf("finished: %s: %s", sig.Name, sig.UUID)
			for _, s := range sig.OnSuccess {
				log.INFO.Printf("sending: %s: %s with payload %v", s.Name, s.RoutingKey, s.Args)
			}
		}

		worker.SetPreTaskHandler(pretaskhandler)
		worker.SetErrorHandler(errorhandler)
		worker.SetPostTaskHandler(posttaskhandler)

		worker.LaunchAsync(errorsChan)
	}
	for {
		select {
		case err := <-errorsChan:
			return err
		}
	}
}

// processTasks are queueing logics
// NOTE it needs to be maintained manually for now
// and server needs to be restarted after updating the tasks
func processTasks() map[string]interface{} {
	airflow := task.NewAirflow()
	argo := task.NewArgo()
	return map[string]interface{}{
		"airflow": airflow.Process,
		"argo":    argo.Process,
	}
}
