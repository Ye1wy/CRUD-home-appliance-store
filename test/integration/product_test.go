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

func (s *TestSuite) TestCreateProdcut() {
	s.CleanTable()
	client := &http.Client{}
	supplierPostUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	supplierData := dto.SupplierRequest{
		Name:        "Narin Inc.",
		PhoneNumber: "66-77-77-13-13",
		Address: &dto.Address{
			Country: "Korea",
			City:    "Seoul",
			Street:  "Dongdaemun",
		},
	}

	supplierPack, err := json.Marshal(&supplierData)
	s.Require().NoError(err)

	supplierPostResp, err := http.Post(supplierPostUrl, "application/json", bytes.NewBuffer(supplierPack))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, supplierPostResp.StatusCode)

	// extraction supplier id
	rawSupplierResp, err := io.ReadAll(supplierPostResp.Body)
	defer supplierPostResp.Body.Close()
	s.Require().NoError(err)

	var supplierResp dto.SupplierResponse // id extract

	err = json.Unmarshal(rawSupplierResp, &supplierResp)
	s.Require().NoError(err)

	path := "../data/bear.png"
	filename := filepath.Base(path)
	buf, err := extractImageData(path, filename)
	s.Require().NoError(err)

	imagePostUrl := fmt.Sprintf("http://%s:%s/api/v1/images", s.cfg.CrudService.Address, s.cfg.CrudService.Port)

	imageReq, err := http.NewRequest(http.MethodPost, imagePostUrl, bytes.NewReader(buf.Bytes()))
	s.Require().NoError(err)

	imageReq.Header.Set("Content-Type", "application/octet-stream")
	imageReq.Header.Set("X-Image-Title", filename)

	imagePostResp, err := client.Do(imageReq)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, imagePostResp.StatusCode)

	// extraction image id
	rawImageResp, err := io.ReadAll(imagePostResp.Body)
	defer imagePostResp.Body.Close()
	s.Require().NoError(err)

	var imageResp dto.ImageResponse // id extract

	err = json.Unmarshal(rawImageResp, &imageResp)
	s.Require().NoError(err)

	productPostUrl := fmt.Sprintf("http://%s:%s/api/v1/products", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	productData := dto.ProductRequest{
		Name:           "Abiba",
		Category:       "Cleaner",
		Price:          120032.23,
		AvailableStock: 999999,
		SupplierId:     supplierResp.Id,
		ImageId:        imageResp.Id,
	}

	productBuf, err := json.Marshal(&productData)
	s.Require().NoError(err)

	productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)

	// extraction image id
	rawProductResp, err := io.ReadAll(productPostResp.Body)
	defer productPostResp.Body.Close()
	s.Require().NoError(err)

	var productResp dto.ProductResponse // id extract

	err = json.Unmarshal(rawProductResp, &productResp)
	s.Require().NoError(err)

	s.Require().NotEmpty(productResp.Id)
	s.Require().Equal(productData.Name, productResp.Name)
	s.Require().Equal(productData.Category, productResp.Category)
	s.Require().Equal(productData.Price, productResp.Price)
	s.Require().Equal(productData.AvailableStock, productResp.AvailableStock)
	s.Require().NotEmpty(productResp.Supplier)
	s.Require().NotEmpty(productResp.Image)
}
