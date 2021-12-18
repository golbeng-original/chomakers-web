package apis

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/golbeng-original/chomakers-web/models"
)

var aboutRepository *models.AboutRepository

func convertResponseAbout(aboutModel *models.AboutModel, aboutHistoryModels []models.AboutHistoryModel) *ResponseAbout {

	responseAbout := ResponseAbout{}

	if aboutModel != nil {
		if aboutModel.ProfileImage != nil {
			responseAbout.ProfileImage = *aboutModel.ProfileImage
		}

		if aboutModel.ProfileName != nil {
			responseAbout.ProfileName = *aboutModel.ProfileName
		}

		if aboutModel.Contact != nil {
			responseAbout.Contact = *aboutModel.Contact
		}

		if aboutModel.IntroduceContent != nil {
			responseAbout.IntroduceContent = *aboutModel.IntroduceContent
		}
	}

	for _, aboutHistoryModel := range aboutHistoryModels {

		aboutHistoryelement := ResponseAboutHistoryElement{
			Id:       aboutHistoryModel.Id,
			Category: aboutHistoryModel.Category,
			Duration: aboutHistoryModel.Duration,
			Content:  aboutHistoryModel.Content,
		}

		responseAbout.Histories = append(responseAbout.Histories, aboutHistoryelement)
	}

	return &responseAbout
}

//func convertResponseAboutHistories(aboutHistoryModel []models.AboutHistoryModel) *

func AboutApis(api *gin.RouterGroup, repositoryConfigure *models.RepositoryConfigure) {

	aboutRepository = repositoryConfigure.AboutRepository

	api.GET("/about", func(c *gin.Context) {

		aboutModel, err := aboutRepository.GetAbout()
		if err != nil {
			errorMessage := fmt.Sprintf("get about error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		aboutHistoryModels, err := aboutRepository.GetHistory()
		if err != nil {
			errorMessage := fmt.Sprintf("get about history error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		responseAbout := convertResponseAbout(aboutModel, aboutHistoryModels)
		responsePresent, err := SuccessResponsePresent(c, responseAbout)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.POST("/about", func(c *gin.Context) {

		var reqAbout RequestUpdateAbout
		c.ShouldBindJSON(&reqAbout)

		complete := false

		var storeImageUrl *string
		if reqAbout.ProfileImage != nil {
			reqSaveImageInfo := models.RequestSaveImageInfo{
				Filename:   reqAbout.ProfileImage.Filename,
				Base64Data: reqAbout.ProfileImage.Data,
			}

			storedImage, err := models.StorageImage("./assets/images", "/images", &reqSaveImageInfo)
			if err != nil {
				errorMessage := fmt.Sprintf("image save error [%v]", err)
				c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
				return
			}

			storeImageUrl = &storedImage.ImageUri

			// 실패하면 지우기...
			defer func(isComplete *bool) {
				if !*isComplete {
					if storedImage != nil {
						os.Remove(storedImage.ImageStorePath)
					}
				}
			}(&complete)
		}

		prevProfileImage, err := aboutRepository.UpdateAbout(storeImageUrl, reqAbout.ProfileName, reqAbout.Contact, reqAbout.IntroduceContent)
		if err != nil {
			errorMessage := fmt.Sprintf("about update error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		if prevProfileImage != nil {
			reletivePath := strings.Replace(*prevProfileImage, "/images", "./assets/images", 1)
			absPath, err := filepath.Abs(reletivePath)
			if err != nil {
				fmt.Println(err)
			} else {
				os.Remove(absPath)
			}
		}

		// 여기까지 오면 성공으로 간주한다.
		complete = true

		aboutModel, err := aboutRepository.GetAbout()
		if err != nil {
			errorMessage := fmt.Sprintf("get about error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		responseAbout := convertResponseAbout(aboutModel, nil)
		responsePresent, err := SuccessResponsePresent(c, responseAbout)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.POST("/about-history", func(c *gin.Context) {

		var reqAboutHistory RequestUpdateAboutHistory
		c.ShouldBindJSON(&reqAboutHistory)

		// 추가되는 내용
		addHistoryInfos := make([]models.AboutHistoryContent, 0)
		for _, responseAddAboutHistory := range reqAboutHistory.AppendHistories {
			addHistoryInfo := models.AboutHistoryContent{
				Category: responseAddAboutHistory.Category,
				Duration: responseAddAboutHistory.Duration,
				Content:  responseAddAboutHistory.Content,
			}

			addHistoryInfos = append(addHistoryInfos, addHistoryInfo)
		}

		// 수정되는 내용
		updateHistoryInfos := make([]models.AboutHistoryIdContent, 0)
		for _, responseUpdateAboutHistory := range reqAboutHistory.UpdateHistories {

			updateHistoryInfo := models.AboutHistoryIdContent{
				Id: responseUpdateAboutHistory.Id,
				AboutHistoryContent: models.AboutHistoryContent{
					Category: responseUpdateAboutHistory.Category,
					Duration: responseUpdateAboutHistory.Duration,
					Content:  responseUpdateAboutHistory.Content,
				},
			}

			updateHistoryInfos = append(updateHistoryInfos, updateHistoryInfo)
		}

		// repositoy 적용
		err := aboutRepository.UpdateAboutHistory(reqAboutHistory.RemoveIds, updateHistoryInfos, addHistoryInfos)
		if err != nil {
			errorMessage := fmt.Sprintf("update aboutHistory error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		aboutHistories, err := aboutRepository.GetHistory()
		if err != nil {
			errorMessage := fmt.Sprintf("get about error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		responseAbout := convertResponseAbout(nil, aboutHistories)
		responsePresent, err := SuccessResponsePresent(c, responseAbout)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})
}
