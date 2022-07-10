package service

import (
	"fmt"
	"log"
	"sut-storage-go/domain/storage/model"
	"sut-storage-go/lib/helper"
	storagepb "sut-storage-go/pb/storage"

	"sut-storage-go/lib/utils"

	"google.golang.org/api/drive/v2"
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

	f := &drive.File{
		Title:    fmt.Sprintf("%s.xlsx", userId),
		MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	p := &drive.ParentReference{Id: s.conf.FolderId}
	f.Parents = []*drive.ParentReference{p}

	res, err := s.driveSrv.Files.
		Insert(f).
		Media(file).
		Do()
	if err != nil {
		log.Fatalln("error when upload file: ", err)
	}

	var fileInfo model.File

	if result := s.H.DB.Where(&model.File{UserId: userId}).First(&fileInfo); result.Error != nil {
		result := s.H.DB.Create(&model.File{
			Id:     res.Id,
			UserId: userId,
		})
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	} else {
		err = s.DeleteFile(ctx, fileInfo.Id)
		if err != nil {
			log.Println(err)
			return err
		}

		result := s.H.DB.Model(&model.File{}).Where(&model.File{UserId: userId}).Update("id", res.Id)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
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
