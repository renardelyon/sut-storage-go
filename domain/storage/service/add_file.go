package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sut-storage-go/lib/helper"
	storagepb "sut-storage-go/pb/storage"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

func (s *Service) AddFile(ctx context.Context, stream storagepb.StorageService_AddFileServer) error {
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

	reader := bytes.NewReader(fileData)

	buf := make([]byte, len(fileData))
	if _, err := reader.Read(buf); err != nil {
		log.Println(err)
		return err
	}

	f := &drive.File{Name: fmt.Sprintf("%s.xlsx", userId), Parents: []string{s.conf.FolderId}}

	res, err := s.driveSrv.Files.
		Create(f).
		Media(reader, googleapi.ContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")).
		Do()
	if err != nil {
		log.Fatalln("error when upload file: ", err)
	}

	// fileId := res.Id
	log.Println(res.Id)
	return nil
}
