package models

import (
	"crypto/sha256"
	base64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type StoredImageInfo struct {
	ImageUri       string
	ImageStorePath string
}

type RequestSaveImageInfo struct {
	Filename   string
	Base64Data string
}

func StorageImages(saveDirectory string, prefixUri string, storageImageInfos []RequestSaveImageInfo) ([]StoredImageInfo, error) {

	storedImages := make([]StoredImageInfo, 0)

	var err error
	for _, storageImageInfo := range storageImageInfos {

		storedImage, err := StorageImage(saveDirectory, prefixUri, &storageImageInfo)
		if err != nil {
			break
		}

		storedImages = append(storedImages, *storedImage)
	}

	if err != nil {

		// err 발생 이전 imageFile들을 제거 한다.
		for _, storedImage := range storedImages {
			os.Remove(storedImage.ImageStorePath)
		}

		return nil, err
	}

	return storedImages, nil
}

func StorageImage(saveDirectory string, prefixUri string, storageImageInfo *RequestSaveImageInfo) (*StoredImageInfo, error) {

	saveDirectory, err := filepath.Abs(saveDirectory)
	if err != nil {
		return nil, err
	}

	bytes, err := base64.StdEncoding.DecodeString(storageImageInfo.Base64Data)
	if err != nil {
		return nil, err
	}

	fileExt := filepath.Ext(storageImageInfo.Filename)

	filename := storageImageInfo.Filename[:len(storageImageInfo.Filename)-len(fileExt)]
	filename = filename + time.Now().String()

	hasher := sha256.New()

	_, err = hasher.Write([]byte(filename))
	if err != nil {
		return nil, err
	}

	hashStr := hasher.Sum(nil)
	stroageFileName := hex.EncodeToString(hashStr)
	stroageFileName = stroageFileName + fileExt

	storageFullPath := fmt.Sprintf("%s/%s", saveDirectory, stroageFileName)
	storageFullPath, err = filepath.Abs(storageFullPath)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(storageFullPath, bytes, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return &StoredImageInfo{
		ImageUri:       fmt.Sprintf("%s/%s", "/images", stroageFileName),
		ImageStorePath: storageFullPath,
	}, nil
}
