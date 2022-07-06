package service

import (
	"sut-storage-go/config"
	"sut-storage-go/lib/pkg/db"

	"google.golang.org/api/drive/v3"
)

type Service struct {
	H        db.Handler
	driveSrv *drive.Service
	conf     *config.Config
}

func NewService(H db.Handler, driveSrv *drive.Service, conf *config.Config) *Service {
	return &Service{
		H,
		driveSrv,
		conf,
	}
}
