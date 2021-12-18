package apis

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"

	"github.com/golbeng-original/chomakers-web/models"
)

var essayRepository *models.EssayRepository

func convertResponseEssayThumbailElement(essayModel *models.EssayThumnailModel) *ResponseEssayThumbnailElement {

	return &ResponseEssayThumbnailElement{
		Id:             essayModel.Id,
		Title:          essayModel.Title,
		ThumbnailImage: essayModel.ThumbnailImage,
	}
}

func convertResponseEssayElement(essayModel *models.EssayModel) *ResponseEssayElement {

	responseImages := make([]ResponseImage, 0)
	for _, imageElement := range essayModel.Images {
		resImageElement := ResponseImage{Id: imageElement.Id, ImageUrl: imageElement.Path}
		responseImages = append(responseImages, resImageElement)
	}

	return &ResponseEssayElement{
		Id:             essayModel.Id,
		Title:          essayModel.Title,
		ThumbnailImage: essayModel.ThumbnailImage,
		Images:         responseImages,
		EssayContent:   essayModel.EssayContent,
	}
}

func EssayApis(api *gin.RouterGroup, repositoryConfigure *models.RepositoryConfigure) {

	essayRepository = repositoryConfigure.EssayRepository

	api.GET("/essay", func(c *gin.Context) {

		allEssayModels, err := essayRepository.GetEssayList()
		if err != nil {
			errorMessage := fmt.Sprintf("get essay list error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		essies := make([]ResponseEssayThumbnailElement, 0)
		for _, essayModel := range allEssayModels {
			essies = append(essies, *convertResponseEssayThumbailElement(&essayModel))
		}

		essayList := &ResponseEssayList{}
		essayList.List = essies

		responsePresent, err := SuccessResponsePresent(c, essayList)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.GET("/essay/:id", func(c *gin.Context) {
		strEssayId := c.Param("id")

		id, err := strconv.Atoi(strEssayId)
		if err != nil {
			errorMessage := fmt.Sprintf("id is wroung (id = %v)", strEssayId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		essayModel, err := essayRepository.FindEssay(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("essay find occur exception [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		resEssayElement := convertResponseEssayElement(essayModel)

		responsePresent, err := SuccessResponsePresent(c, resEssayElement)
		if err != nil {
			errorMessage := fmt.Sprintf("get one essay error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.POST("/essay", func(c *gin.Context) {

		var reqCreateEssay RequestCreateEssay
		c.ShouldBindJSON(&reqCreateEssay)

		requestThumbnailSaveImageInfo := models.RequestSaveImageInfo{
			Filename:   reqCreateEssay.ThumbnailImage.Filename,
			Base64Data: reqCreateEssay.ThumbnailImage.Data,
		}

		storedThumbnailImage, err := models.StorageImage("./assets/images", "/images", &requestThumbnailSaveImageInfo)
		if err != nil {
			errorMessage := fmt.Sprintf("thumbnail image save error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		requestSaveImageInfos := make([]models.RequestSaveImageInfo, 0)
		for _, reqImage := range reqCreateEssay.Images {
			requestSaveImageInfos = append(requestSaveImageInfos, models.RequestSaveImageInfo{Filename: reqImage.Filename, Base64Data: reqImage.Data})
		}

		storedImages, err := models.StorageImages("./assets/images", "/images", requestSaveImageInfos)
		if err != nil {
			errorMessage := fmt.Sprintf("image save error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		images := funk.Map(storedImages, func(e models.StoredImageInfo) string {
			return e.ImageUri
		})

		insertId, err := essayRepository.AddEssay(
			reqCreateEssay.Title,
			storedThumbnailImage.ImageUri,
			reqCreateEssay.EssayContent,
			images.([]string),
		)

		if err != nil {
			errorMessage := fmt.Sprintf("AddPotofolio error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		essayModel, err := essayRepository.FindEssay(insertId)
		if err != nil {
			errorMessage := fmt.Sprintf("AddEssay after error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		responseEssayElement := convertResponseEssayElement(essayModel)
		responsePresent, err := SuccessResponsePresent(c, responseEssayElement)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.PUT("essay/:id", func(c *gin.Context) {

		complete := false

		strPotofolioId := c.Param("id")
		id, err := strconv.Atoi(strPotofolioId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"result": "failed",
				"error":  fmt.Sprintf("id is wroung (id = %s)", strPotofolioId),
			})
		}

		prevEssayModel, err := essayRepository.FindEssay(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("essayId = %d [err = %s]", id, err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		var requestUpdateEssay RequestUpdateEssay
		c.ShouldBindJSON(&requestUpdateEssay)

		// thumbnail image update
		var prevThumbnailImagePath *string
		var sotredThumbnailImagePath *models.StoredImageInfo
		var storedthumbnailUrl *string
		if requestUpdateEssay.NewThumbnail != nil {

			prevThumbnailImagePath = new(string)
			*prevThumbnailImagePath = prevEssayModel.GetThumbnailImagePath("./assets/images", "/image")

			requestSaveImageInfo := models.RequestSaveImageInfo{
				Filename:   requestUpdateEssay.NewThumbnail.Filename,
				Base64Data: requestUpdateEssay.NewThumbnail.Data,
			}

			sotredThumbnailImagePath, err = models.StorageImage("./assets/images", "/images", &requestSaveImageInfo)
			if err != nil {
				errorMessage := fmt.Sprintf("thumbnailImage save error [%v]", err)
				c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
				return
			}

			storedthumbnailUrl = &sotredThumbnailImagePath.ImageUri
		}

		// 성공/실패 여부에 따른 Thumbnail 이미지 처리
		defer func(isComplete *bool) {
			if *isComplete {
				if prevThumbnailImagePath != nil {
					os.Remove(*prevThumbnailImagePath)
				}
			} else {
				if sotredThumbnailImagePath != nil {
					os.Remove(sotredThumbnailImagePath.ImageStorePath)
				}
			}
		}(&complete)

		// images update
		var storedImages []models.StoredImageInfo
		var images interface{}
		if requestUpdateEssay.AddImages != nil {

			requestSaveImageInfos := make([]models.RequestSaveImageInfo, 0)
			for _, reqImage := range requestUpdateEssay.AddImages {
				requestSaveImageInfos = append(requestSaveImageInfos, models.RequestSaveImageInfo{Filename: reqImage.Filename, Base64Data: reqImage.Data})
			}

			storedImages, err = models.StorageImages("./assets/images", "/images", requestSaveImageInfos)
			if err != nil {
				errorMessage := fmt.Sprintf("image save error [%v]", err)
				c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
				return
			}

			images = funk.Map(storedImages, func(e models.StoredImageInfo) string {
				return e.ImageUri
			})
		}

		// 성공/실패 여부에 따른 Thumbnail 이미지 처리
		defer func(isComplete *bool) {
			if !*isComplete {
				for _, sotredImage := range storedImages {
					os.Remove(sotredImage.ImageStorePath)
				}
			}

		}(&complete)

		removeImages, err := essayRepository.UpdateEssay(int64(id),
			requestUpdateEssay.Title,
			storedthumbnailUrl,
			requestUpdateEssay.EssayContent,
			requestUpdateEssay.RemoveImageIds,
			images.([]string))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"result": "failed",
				"error":  fmt.Sprintf("potofolioId = %d repository.UpdateEssay [err = %s]", id, err),
			})
		}

		// 파일 지우기
		for _, removeImage := range removeImages {
			reletivePath := strings.Replace(removeImage, "/images", "./assets/images", 1)
			absPath, err := filepath.Abs(reletivePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			os.Remove(absPath)
		}

		// 여기까지 오면 성공으로 간주한다.
		complete = true

		essayModel, err := essayRepository.FindEssay(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("Update essay after error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		responseEssayElment := convertResponseEssayElement(essayModel)
		responsePresent, err := SuccessResponsePresent(c, responseEssayElment)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.DELETE("/essay/:id", func(c *gin.Context) {
		strPotofolioId := c.Param("id")

		id, err := strconv.Atoi(strPotofolioId)
		if err != nil {
			errorMessage := fmt.Sprintf("id is wroung (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		_, err = essayRepository.FindEssay(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("essay not found (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		err = essayRepository.RemoveEssay(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("essayRepository.RemoveEssay (id = %s) [%v]", strPotofolioId, err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		responsePresent, err := SuccessResponsePresent(c, nil)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

}
