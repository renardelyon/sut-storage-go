package service

import (
	"fmt"
	"log"
	"sut-storage-go/lib/helper"
	storagepb "sut-storage-go/pb/storage"

	"sut-storage-go/lib/utils"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

func (s *Service) AddFile(stream storagepb.StorageService_AddFileServer) error {
	ctx := stream.Context()
	select {
	case <-ctx.Done():
		return helper.ContextError(ctx)
	default:
	}

	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}

	if req.GetInfo() == nil {
		log.Println(err)
		return err
	}

	userId := req.GetInfo().UserId

	fileData, err := s.ReceiveStreamFile(ctx, stream)
	if err != nil {
		log.Println(err)
		return err
	}

	tmpFilePath := fmt.Sprintf("/tmp/%s.xlsx", userId)

	tmpfCreator := utils.TempFileCreator{
		TempFilePath: tmpFilePath,
	}

	file, err := tmpfCreator.CreateTempFile(fileData)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	defer tmpfCreator.DeleteTempFile()

	f := &drive.File{Name: fmt.Sprintf("%s.xlsx", userId), Parents: []string{s.conf.FolderId}}

	res, err := s.driveSrv.Files.
		Create(f).
		Media(file, googleapi.ContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")).
		Do()
	if err != nil {
		log.Fatalln("error when upload file: ", err)
	}

	err = stream.SendAndClose(&storagepb.UploadResponse{
		Id: res.Id,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
