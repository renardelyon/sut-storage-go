package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sut-storage-go/lib/helper"
	"sut-storage-go/lib/pkg/gdrive"
	storagepb "sut-storage-go/pb/storage"

	"google.golang.org/api/drive/v2"
)

func (s *Service) ReceiveStreamFile(ctx context.Context, stream storagepb.StorageService_AddFileServer) ([]byte, error) {
	const maxFileSize = 1 << 30 // 1GB

	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		select {
		case <-ctx.Done():
			return nil, helper.ContextError(ctx)
		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more stream data")
			break
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		fileSize += size
		if fileSize > maxFileSize {
			log.Println("File is larger than 1 GB")
			return nil, errors.New("file is larger than 1 GB")
		}

		_, err = fileData.Write(chunk)
		if err != nil {
			return nil, err
		}
	}

	return fileData.Bytes(), nil
}

func (s *Service) DeleteFile(ctx context.Context, fileId string) error {
	select {
	case <-ctx.Done():
		return helper.ContextError(ctx)
	default:
	}

	err := s.driveSrv.Children.Delete(s.conf.FolderId, fileId).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) handleTokenRegeneration(userId string, file *os.File) *drive.File {
	handler := gdrive.DriveHandler{
		Config: s.conf,
	}

	tok, err := handler.RegenerateToken()
	if err != nil {
		log.Fatalln(err)
	}
	gdrive.SaveToken(s.conf.TokenPath, tok)

	srv, err := handler.NewDriveService()
	if err != nil {
		log.Fatalln(err)
	}

	f := &drive.File{
		Title:    fmt.Sprintf("%s.xlsx", userId),
		MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	p := &drive.ParentReference{Id: s.conf.FolderId}
	f.Parents = []*drive.ParentReference{p}

	res, err := srv.Files.
		Insert(f).
		Media(file).
		Do()
	if err != nil {
		log.Fatalln("error when upload file: ", err)
	}

	return res
}
