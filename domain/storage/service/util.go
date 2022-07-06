package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"sut-storage-go/lib/helper"
	storagepb "sut-storage-go/pb/storage"
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
