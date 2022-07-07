package utils

import (
	"log"
	"os"
)

type TempFileCreator struct {
	TempFilePath string
}

func (temp *TempFileCreator) CreateTempFile(fileData []byte) (*os.File, error) {
	fileCreate, err := os.Create(temp.TempFilePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fileCreate.Close()

	_, err = fileCreate.Write(fileData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open(temp.TempFilePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return file, nil
}

func (temp *TempFileCreator) DeleteTempFile() error {
	err := os.Remove(temp.TempFilePath)
	if err != nil {
		return err
	}

	return nil
}
