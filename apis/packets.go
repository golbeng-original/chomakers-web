package apis

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Save Image Common
type RequestSaveImage struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

type ResponseImage struct {
	Id       int64  `json:"id"`
	ImageUrl string `json:"image"`
}

// Login
type RequestLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	LoginResult int `json:"result"` // 0: Sucess, 1: wroung username, 2; wroung password, 3: token generate fail
}

// Potoflio Get
type ResponsePotofolioElement struct {
	Id     int64           `json:"id"`
	Title  string          `json:"title"`
	Images []ResponseImage `json:"images"`
}

type ResponsePotofolioList struct {
	List []ResponsePotofolioElement `json:"list"`
}

// Potofolio Update (PUT)
type RequestUpdatePotofolio struct {
	Title          *string            `json:"title"`
	RemoveImageIds []int64            `json:"remove_images"`
	AddImages      []RequestSaveImage `json:"add_images"`
}

// Potofolio New (Post)
type RequestCreatePotofolio struct {
	Title  string             `json:"title"`
	Images []RequestSaveImage `json:"images"`
}

// Essay Thumbnail
type ResponseEssayThumbnailElement struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	ThumbnailImage string `json:"thumbnail"`
}

type ResponseEssayList struct {
	List []ResponseEssayThumbnailElement `json:"list"`
}

// Essay
type ResponseEssayElement struct {
	Id             int64           `json:"id"`
	Title          string          `json:"title"`
	ThumbnailImage string          `json:"thumbmail"`
	Images         []ResponseImage `json:"images"`
	EssayContent   string          `json:"essay_content"`
}

// Essay New (Post)
type RequestCreateEssay struct {
	Title          string             `json:"title"`
	ThumbnailImage RequestSaveImage   `json:"thumbmail"`
	Images         []RequestSaveImage `json:"images"`
	EssayContent   string             `json:"essay_content"`
}

// Essay Update (PUT)
type RequestUpdateEssay struct {
	Title          *string            `json:"title"`
	NewThumbnail   *RequestSaveImage  `json:"thumbnail"`
	RemoveImageIds []int64            `json:"remove_images"`
	AddImages      []RequestSaveImage `json:"add_images"`
	EssayContent   *string            `json:"essay_content"`
}

// About History
type ResponseAboutHistoryElement struct {
	Id       int64  `json:"id"`
	Category string `json:"category"`
	Duration string `json:"duration"`
	Content  string `json:"content"`
}

// About
type ResponseAbout struct {
	ProfileImage     string `json:"profile_image"`
	ProfileName      string `json:"profile_name"`
	Contact          string `json:"contact"`
	IntroduceContent string `json:"introduce_content"`

	Histories []ResponseAboutHistoryElement `json:"history_list"`
}

// Update About
type RequestUpdateAbout struct {
	ProfileImage     *RequestSaveImage `json:"profile_image"`
	ProfileName      *string           `json:"profile_name"`
	Contact          *string           `json:"contact"`
	IntroduceContent *string           `json:"introduce_content"`
}

// Update About History
type RequesAboutHistoryElement struct {
	Category string `json:"category"`
	Duration string `json:"duration"`
	Content  string `json:"content"`
}

type RequestUpdateAboutHistoryElement struct {
	Id int64 `json:"id"`
	RequesAboutHistoryElement
}

type RequestUpdateAboutHistory struct {
	RemoveIds       []int64                            `json:"remove_id_list"`
	AppendHistories []RequesAboutHistoryElement        `json:"append_history_list"`
	UpdateHistories []RequestUpdateAboutHistoryElement `json:"update_history_list"`
}

//
type ResponsePresent struct {
	Result string `json:"result"`
	Header string `json:"header"`
	Error  string `json:"error"`
	Data   string `json:"data"`
}

func SuccessResponsePresent(c *gin.Context, data interface{}) (*ResponsePresent, error) {

	headers := c.Writer.Header()

	fmt.Println(len(headers))

	responseHeaders := make([]string, 0)

	for key, header := range headers {
		responseHeader := key

		for _, value := range header {

			responseHeader = fmt.Sprintf("%s:%s", key, value)
			responseHeaders = append(responseHeaders, responseHeader)
		}
	}

	headerJsonBytes, _ := json.Marshal(responseHeaders)
	headerStr := string(headerJsonBytes)

	bodyStr := "'"
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		bodyStr = string(bytes)
	}

	response := ResponsePresent{Result: "success", Data: bodyStr, Header: headerStr}
	return &response, nil
}

func FailedResponsePreset(err string) *ResponsePresent {
	response := ResponsePresent{Result: "failed", Error: err}
	return &response
}
