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

var potofolioRepository *models.PotofolioRepository

func convertResponsePotofolioElement(potofolioModel *models.PotofolioModel) *ResponsePotofolioElement {

	responseImages := make([]ResponseImage, 0)

	for _, imageElement := range potofolioModel.Images {
		resImageElement := ResponseImage{Id: imageElement.Id, ImageUrl: imageElement.Path}
		responseImages = append(responseImages, resImageElement)
	}

	return &ResponsePotofolioElement{
		Id:     potofolioModel.Id,
		Title:  potofolioModel.Title,
		Images: responseImages,
	}
}

func PotofolioApis(api *gin.RouterGroup, repositoryConfigure *models.RepositoryConfigure) {

	potofolioRepository = repositoryConfigure.PotofolioRepository

	api.GET("/potofolio/:id", func(c *gin.Context) {
		strPotofolioId := c.Param("id")

		id, err := strconv.Atoi(strPotofolioId)
		if err != nil {
			errorMessage := fmt.Sprintf("id is wroung (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		potofolioModel, err := potofolioRepository.FindPotofolio(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("potofolio find occur exception [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		resPotofolioElement := convertResponsePotofolioElement(potofolioModel)

		responsePresent, err := SuccessResponsePresent(c, resPotofolioElement)
		if err != nil {
			errorMessage := fmt.Sprintf("get one potofolio error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.GET("/potofolio", func(c *gin.Context) {

		allPotofolioModels, err := potofolioRepository.GetPotofolioList()
		if err != nil {
			errorMessage := fmt.Sprintf("get potofolio list error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		potofolios := make([]ResponsePotofolioElement, 0)
		for _, potofolioModel := range allPotofolioModels {
			potofolios = append(potofolios, *convertResponsePotofolioElement(&potofolioModel))
		}

		potofolioList := &ResponsePotofolioList{}
		potofolioList.List = potofolios

		responsePresent, err := SuccessResponsePresent(c, potofolioList)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	// 생성
	api.POST("/potofolio", func(c *gin.Context) {

		var reqCreatePotofolio RequestCreatePotofolio
		c.ShouldBindJSON(&reqCreatePotofolio)

		requestSaveImageInfos := make([]models.RequestSaveImageInfo, 0)
		for _, reqImage := range reqCreatePotofolio.Images {
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

		insertId, err := potofolioRepository.AddPotofolio(reqCreatePotofolio.Title, images.([]string))
		if err != nil {
			errorMessage := fmt.Sprintf("AddPotofolio error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		potofolioModel, err := potofolioRepository.FindPotofolio(insertId)
		if err != nil {
			errorMessage := fmt.Sprintf("AddPotofolio after error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		responsePotofolioElment := convertResponsePotofolioElement(potofolioModel)
		responsePresent, err := SuccessResponsePresent(c, responsePotofolioElment)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	// 수정
	api.PUT("/potofolio/:id", func(c *gin.Context) {

		complete := false

		strPotofolioId := c.Param("id")

		id, err := strconv.Atoi(strPotofolioId)
		if err != nil {
			errorMessage := fmt.Sprintf("id is wroung (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		_, err = potofolioRepository.FindPotofolio(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("potofolioId = %d [err = %s]", id, err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		var reqUpdatePotofolio RequestUpdatePotofolio
		c.ShouldBindJSON(&reqUpdatePotofolio)

		var storedImages []models.StoredImageInfo
		var images interface{}
		if reqUpdatePotofolio.AddImages != nil {

			requestSaveImageInfos := make([]models.RequestSaveImageInfo, 0)
			for _, reqImage := range reqUpdatePotofolio.AddImages {
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

		defer func(isComplete *bool) {
			if !*isComplete {
				for _, sotredImage := range storedImages {
					os.Remove(sotredImage.ImageStorePath)
				}
			}

		}(&complete)

		removeImages, err := potofolioRepository.UpdatePotofolio(int64(id), reqUpdatePotofolio.Title, reqUpdatePotofolio.RemoveImageIds, images.([]string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"result": "failed",
				"error":  fmt.Sprintf("potofolioId = %d repository.UpdatePotofolio [err = %s]", id, err),
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

		potofolioModel, err := potofolioRepository.FindPotofolio(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("Update Potofolio after error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		responsePotofolioElment := convertResponsePotofolioElement(potofolioModel)
		responsePresent, err := SuccessResponsePresent(c, responsePotofolioElment)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	// 제거
	api.DELETE("/potofolio/:id", func(c *gin.Context) {

		strPotofolioId := c.Param("id")

		id, err := strconv.Atoi(strPotofolioId)
		if err != nil {
			errorMessage := fmt.Sprintf("id is wroung (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		_, err = potofolioRepository.FindPotofolio(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("potofolio not found (id = %s)", strPotofolioId)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
		}

		err = potofolioRepository.RemovePotofolio(int64(id))
		if err != nil {
			errorMessage := fmt.Sprintf("potofolioRepository.RemovePotofolio (id = %s) [%v]", strPotofolioId, err)
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
