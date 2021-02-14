package main

import (
	"flag"
	"kudosbot/kudos"
	"kudosbot/pkg/config"
	"kudosbot/pkg/database"
	"kudosbot/repository"

	"go.uber.org/fx"
)

func main() {

	configFile := flag.String("config", "config/config.yaml.local", "config file")
	flag.Parse()

	app := fx.New(
		fx.Provide(
			func() *config.Config { return config.NewConfig(*configFile) },
			database.NewDatabase,
			kudos.NewService,
			repository.NewKudosRepository,
		),
		fx.Invoke(Start))
	app.Run()
	<-app.Done()
}

func Start(
	service *kudos.Service,
) {
	go service.Run()
}
