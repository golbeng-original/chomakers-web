package apis

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponsePotofolioList(t *testing.T) {

	var potofolios []ResponsePotofolioElement

	responseImage_1 := ResponseImage{Id: 0, ImageUrl: "url1"}
	responseImage_2 := ResponseImage{Id: 0, ImageUrl: "url2"}
	responseImage_3 := ResponseImage{Id: 0, ImageUrl: "url3"}
	responseImage_4 := ResponseImage{Id: 0, ImageUrl: "url4"}

	potofolios = append(potofolios, ResponsePotofolioElement{Id: 0, Title: "test_title0", Images: []ResponseImage{responseImage_1, responseImage_2}})
	potofolios = append(potofolios, ResponsePotofolioElement{Id: 1, Title: "test_title1", Images: []ResponseImage{responseImage_3, responseImage_4}})

	jsonString, err := json.Marshal(&potofolios)
	assert.Nil(t, err)

	fmt.Println(string(jsonString))
}

func TestRequestPtotolioUpdateAll(t *testing.T) {

	jsonString := `{
			"title": "test_title",
			"remove_images" : [1,2,3,4],
			"add_images" : [
				{
					"filename" : "test_file_name_1.png",
					"data" : "data!!!"
				},
				{
					"filename" : "test_file_name_2.png",
					"data" : "data2!!!"
				}
			]
		}`

	var updateRequest RequestUpdatePotofolio

	err := json.Unmarshal([]byte(jsonString), &updateRequest)
	assert.Nil(t, err)

	assert.Equal(t, *updateRequest.Title, "test_title")
	assert.Equal(t, updateRequest.RemoveImageIds, []int64{1, 2, 3, 4})

	assert.Equal(t, len(updateRequest.AddImages), 2)

	assert.Equal(t, (updateRequest.AddImages)[0].Filename, "test_file_name_1.png")
	assert.Equal(t, (updateRequest.AddImages)[0].Data, "data!!!")
}

func TestRequestPtotolioUpdatePartial(t *testing.T) {

	jsonString := `{
			"add_images" : [
				{
					"filename" : "test_file_name_1.png",
					"data" : "data!!!"
				},
				{
					"filename" : "test_file_name_2.png",
					"data" : "data2!!!"
				}
			]
		}`

	var updateRequest RequestUpdatePotofolio

	err := json.Unmarshal([]byte(jsonString), &updateRequest)
	assert.Nil(t, err)

	assert.Nil(t, updateRequest.Title)
	assert.Nil(t, updateRequest.RemoveImageIds)

	assert.Equal(t, len(updateRequest.AddImages), 2)

	assert.Equal(t, (updateRequest.AddImages)[0].Filename, "test_file_name_1.png")
	assert.Equal(t, (updateRequest.AddImages)[0].Data, "data!!!")
}
