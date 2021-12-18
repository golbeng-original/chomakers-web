package models

import (
	base64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageImage(t *testing.T) {

	filename := "a0.jpg"

	originalImagePath := fmt.Sprintf("../assets/test_images/%s", filename)

	bytes, err := ioutil.ReadFile(originalImagePath)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	base64Bytes := base64.StdEncoding.EncodeToString(bytes)
	assert.NotNil(t, base64Bytes)

	reqSaveImageInfo := RequestSaveImageInfo{
		Filename:   filename,
		Base64Data: base64Bytes,
	}

	storageFileName, err := StorageImage("../assets/test_images", filename, &reqSaveImageInfo)
	assert.Nil(t, err)
	assert.NotNil(t, storageFileName)

	_, err = os.Stat(storageFileName.ImageStorePath)
	assert.Nil(t, err)

	os.Remove(storageFileName.ImageStorePath)

	//fmt.Printf("storageFileName = %s\n", storageFileName)
}
