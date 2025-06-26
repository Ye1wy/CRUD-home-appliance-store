package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

func (s *TestSuite) TestCreateImage() {
	s.CleanTable()
	buf, err := extractImageData("../data/cat.png", "cat.png")
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf.Bytes()))
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

	var extracted []dto.Image
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 1)
	s.Require().Equal("cat", extracted[0].Title)
}

func (s *TestSuite) TestCreateImageEmptyData() {
	s.CleanTable()
	url := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	req, err := http.NewRequest("POST", url, nil)
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
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf.Bytes()))
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
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf.Bytes()))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	postResp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, postResp.StatusCode)
}

func (s *TestSuite) TestGetAll() {
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
		req, err := http.NewRequest("POST", url, bytes.NewReader(buf.Bytes()))
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

	var extracted []dto.Image
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 10)
}

func (s *TestSuite) TestGetAllWithLimit() {
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
		req, err := http.NewRequest("POST", postUrl, bytes.NewReader(buf.Bytes()))
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

	var extracted []dto.Image
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 17)
}

func (s *TestSuite) TestGetAllWithOffset() {
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
		req, err := http.NewRequest("POST", postUrl, bytes.NewReader(buf.Bytes()))
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

	var extracted []dto.Image
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 7)
}

func (s *TestSuite) TestGetAllWithLimitAndOffset() {
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
		req, err := http.NewRequest("POST", postUrl, bytes.NewReader(buf.Bytes()))
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

	var extracted []dto.Image
	err = json.Unmarshal(rawData, &extracted)
	s.Require().NoError(err)
	s.Require().Len(extracted, 12)
}

func (s *TestSuite) TestGetId() {}
