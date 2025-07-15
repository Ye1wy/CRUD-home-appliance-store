package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *TestSuite) TestCreateImage() {
	s.CleanTable()
	buf, err := extractImageData("../data/cat.png", "cat.png")
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf.Bytes()))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Image-Title", "cat")
	client := &http.Client{}
	postResp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, postResp.StatusCode)

	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var extracted []dto.ImageResponse
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 1)
	s.Require().Equal("cat", extracted[0].Title)
}

func (s *TestSuite) TestCreateImageEmptyData() {
	s.CleanTable()
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Image-Title", "cat")
	client := &http.Client{}
	postResp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, postResp.StatusCode)
}

func (s *TestSuite) TestCreateImageWithEmptyTitle() {
	s.CleanTable()
	buf, err := extractImageData("../data/cat.png", "cat.png")
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf.Bytes()))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{}
	postResp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, postResp.StatusCode)
}

func (s *TestSuite) TestCreateImageAnotherContentType() {
	s.CleanTable()
	buf, err := extractImageData("../data/cat.png", "cat.png")
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf.Bytes()))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	postResp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, postResp.StatusCode)
}

func (s *TestSuite) TestGetAllImage() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	getResp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var extracted []dto.ImageResponse
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 10)
}

func (s *TestSuite) TestGetAllImageWithLimit() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	getUrl := postUrl + "?limit=17"
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var extracted []dto.ImageResponse
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 17)
}

func (s *TestSuite) TestGetAllImageWithOffset() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	getUrl := postUrl + "?offset=10"
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var extracted []dto.ImageResponse
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 7)
}

func (s *TestSuite) TestGetAllImageWithLimitAndOffset() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	getUrl := postUrl + "?limit=20&offset=5"
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var extracted []dto.ImageResponse
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 12)
}

func (s *TestSuite) TestGetByIdImage() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	sql := `SELECT * FROM image`
	rows, err := s.db.Query(context.Background(), sql)
	s.Require().NoError(err)

	var imageData []domain.Image

	for rows.Next() {
		var temp domain.Image

		err := rows.Scan(&temp.Id, &temp.Title, &temp.Data)
		s.Require().NoError(err)

		imageData = append(imageData, temp)
	}

	checkInstance := imageData[len(imageData)/2]

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, checkInstance.Id.String())
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)
	defer getResp.Body.Close()

	buf, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)
	var takedData dto.ImageResponse
	err = json.Unmarshal(buf, &takedData)
	s.Require().NoError(err)

	s.Require().Equal(checkInstance.Title, takedData.Title)
	s.Require().Equal(checkInstance.Data, takedData.Image)
}

func (s *TestSuite) TestGetByIdImageNonValidId() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, "ads1231asda1")
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, getResp.StatusCode)
}

func (s *TestSuite) TestGetByIdImageNonContainedId() {
	s.CleanTable()
	allData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
		"../data/oni_girl.jpg",
		"../data/mem-4.jpg",
		"../data/mem-6.jpg",
		"../data/mem-2.jpg",
		"../data/miyabi.jpg",
		"../data/mem-7.png",
		"../data/miku.jpg",
		"../data/tat.jpg",
		"../data/mem-5.jpg",
		"../data/feelsbadman.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range allData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	randomUUID, err := uuid.NewRandom()
	s.Require().NoError(err)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, randomUUID.String())
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *TestSuite) TestUpdateImage() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	sqlQuery := `SELECT * FROM image`
	rows, err := s.db.Query(context.Background(), sqlQuery)
	s.Require().NoError(err)

	var extractedData []domain.Image

	for rows.Next() {
		var temp domain.Image
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Data)
		s.Require().NoError(err)

		extractedData = append(extractedData, temp)
	}

	randomId := rand.IntN(len(extractedData))
	updateData := extractedData[randomId]

	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, updateData.Id.String())

	updateImage := "../data/vash.png"
	patchBody, err := extractImageData(updateImage, filepath.Base(updateImage))
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPatch, idxUrl, patchBody)
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/octet-stream")
	httpClient := &http.Client{}
	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	sqlCheck := `SELECT COUNT(id) FROM image`
	var count int
	err = s.db.QueryRow(context.Background(), sqlCheck).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(len(extractedData), count)

	getResp, err := http.Get(idxUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)
	defer getResp.Body.Close()

	buf, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var receivedImage domain.Image
	err = json.Unmarshal(buf, &receivedImage)
	s.Require().NoError(err)
	s.Require().Equal(updateData.Title, receivedImage.Title)
	expectedImageHash := hashBytes(patchBody.Bytes())
	takedImageHash := hashBytes(receivedImage.Data)
	s.Require().Equal(expectedImageHash, takedImageHash)
}

func (s *TestSuite) TestUpdateImageNonContainedId() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	nonContainedId, err := uuid.NewRandom()
	s.Require().NoError(err)
	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, nonContainedId.String())

	updateImage := "../data/vash.png"
	patchBody, err := extractImageData(updateImage, filepath.Base(updateImage))
	s.Require().NoError(err)

	patchReq, err := http.NewRequest(http.MethodPatch, idxUrl, patchBody)
	patchReq.Header.Set("Content-Type", "application/octet-stream")
	s.Require().NoError(err)
	client := http.Client{}
	resp, err := client.Do(patchReq)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *TestSuite) TestUpdateImageInvalidId() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, "aboba1337228")
	patchReq, err := http.NewRequest(http.MethodPatch, idxUrl, nil)
	s.Require().NoError(err)
	client := http.Client{}
	resp, err := client.Do(patchReq)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestSuite) TestUpdateImageWithChangeTitle() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	sqlQuery := `SELECT * FROM image`
	rows, err := s.db.Query(context.Background(), sqlQuery)
	s.Require().NoError(err)

	var extractedData []domain.Image

	for rows.Next() {
		var temp domain.Image
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Data)
		s.Require().NoError(err)

		extractedData = append(extractedData, temp)
	}

	randomId := rand.IntN(len(extractedData))
	updateData := extractedData[randomId]

	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, updateData.Id.String())

	updateImage := "../data/vash.png"
	patchBody, err := extractImageData(updateImage, filepath.Base(updateImage))
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPatch, idxUrl, patchBody)
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/octet-stream")
	req.Header.Set("X-Image-Title", "Aboba")
	httpClient := &http.Client{}
	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	sqlCheckCount := `SELECT COUNT(id) FROM image`
	var count int
	err = s.db.QueryRow(context.Background(), sqlCheckCount).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(len(extractedData), count)

	getResp, err := http.Get(idxUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)
	defer getResp.Body.Close()

	buf, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var receivedImage domain.Image
	err = json.Unmarshal(buf, &receivedImage)
	s.Require().NoError(err)
	expectedImageHash := hashBytes(patchBody.Bytes())
	takedImageHash := hashBytes(receivedImage.Data)
	s.Require().Equal(expectedImageHash, takedImageHash)

	sqlCheckContent := `SELECT * FROM image WHERE id=@checkId`
	arg := pgx.NamedArgs{
		"checkId": updateData.Id.String(),
	}

	row := s.db.QueryRow(context.Background(), sqlCheckContent, arg)
	var changedData domain.Image
	err = row.Scan(
		&changedData.Id,
		&changedData.Title,
		&changedData.Data,
	)
	s.Require().NoError(err)
	s.Require().Equal("Aboba", changedData.Title)
}

func (s *TestSuite) TestDeleteImage() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	sqlQuery := `SELECT * FROM image`
	rows, err := s.db.Query(context.Background(), sqlQuery)
	s.Require().NoError(err)

	var extractedData []domain.Image

	for rows.Next() {
		var temp domain.Image
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Data)
		s.Require().NoError(err)

		extractedData = append(extractedData, temp)
	}

	randomId := rand.IntN(len(extractedData))
	deleteData := extractedData[randomId]

	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, deleteData.Id.String())
	req, err := http.NewRequest(http.MethodDelete, idxUrl, nil)
	s.Require().NoError(err)
	httpClient := &http.Client{}
	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, patchResp.StatusCode)

	sqlCheckCount := `SELECT COUNT(id) FROM image`
	var count int
	err = s.db.QueryRow(context.Background(), sqlCheckCount).Scan(&count)
	s.Require().NoError(err)
	expectedCount := len(extractedData) - 1
	s.Require().Equal(expectedCount, count)
}

func (s *TestSuite) TestDeleteImageNonContainedId() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	randomUUID, err := uuid.NewRandom()
	s.Require().NoError(err)
	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, randomUUID.String())
	req, err := http.NewRequest(http.MethodDelete, idxUrl, nil)
	s.Require().NoError(err)
	httpClient := &http.Client{}
	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, patchResp.StatusCode)
}

func (s *TestSuite) TestDeleteImageInvalidId() {
	s.CleanTable()
	loadedData := []string{
		"../data/cat.png",
		"../data/vash.png",
		"../data/kosmodes.png",
		"../data/mem-3.jpg",
		"../data/gen.png",
		"../data/mem-1.jpg",
		"../data/bear.png",
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	for _, path := range loadedData {
		filename := filepath.Base(path)
		buf, err := extractImageData(path, filename)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(buf.Bytes()))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Image-Title", filename)
		client := &http.Client{}
		postResp, err := client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, postResp.StatusCode)
	}

	idxUrl := fmt.Sprintf("http://%s:%s/api/v1/images/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, "abuba12")
	req, err := http.NewRequest(http.MethodDelete, idxUrl, nil)
	s.Require().NoError(err)
	httpClient := &http.Client{}
	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, patchResp.StatusCode)
}
