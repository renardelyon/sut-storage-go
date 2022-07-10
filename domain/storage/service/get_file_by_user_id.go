package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"sut-storage-go/domain/storage/model"
	"sut-storage-go/lib/utils"
	storagepb "sut-storage-go/pb/storage"

	"google.golang.org/api/drive/v2"
)

func (s *Service) GetFileByUserId(req *storagepb.GetFileByUserIdRequest, stream storagepb.StorageService_GetFileByUserIdServer) error {
	var fileInfo model.File
	if result := s.H.DB.Where(&model.File{UserId: req.UserId}).First(&fileInfo); result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	f := &drive.File{
		DownloadUrl: fmt.Sprintf("https://docs.google.com/uc?export=download&id=%s", fileInfo.Id),
	}

	resFile, err := utils.DownloadFile(f)
	if err != nil {
		log.Println(err)
		return err
	}

	bytesReader := bytes.NewReader(resFile)
	reader := bufio.NewReader(bytesReader)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return err
		}

		req := &storagepb.GetFileByUserIdResponse{
			ChunkData: buffer[:n],
		}

		err = stream.Send(req)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// tmp, _ := os.Create("coba.xlsx")
	// defer tmp.Close()
	// fileSize, err := tmp.Write(resFile)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// log.Println("fileSize", fileSize)
	// tmp.Sync()

	return nil
}
