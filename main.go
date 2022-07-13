package main

import (
	"log"
	"sut-storage-go/application"
	"sut-storage-go/config"
	"sut-storage-go/domain/storage/service"
	storagepb "sut-storage-go/pb/storage"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config: ", err.Error())
	}

	app, err := application.Setup(&c)
	if err != nil {
		log.Fatalln("Failed at application setup: ", err.Error())
	}

	s := service.NewService(app.DbClients, app.DriveHandler, &c)

	storagepb.RegisterStorageServiceServer(app.GrpcServer, s)

	err = app.Run(&c)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
