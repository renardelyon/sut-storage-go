package application

import "sut-storage-go/config"

func (app *Application) Run(cfg *config.Config) error {
	return grpcRun(cfg)(app)
}
