package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
)

func productResponseToRequest(product dto.ProductResponse) dto.ProductRequest {
	out := dto.ProductRequest{
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		SupplierId:     product.Supplier.Id,
		ImageId:        product.Image.Id,
	}

	return out
}

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

func (s *TestSuite) TestCreateProdcutWithoutImageIdAndSupplierId() {
	s.CleanTable()

	productPostUrl := fmt.Sprintf("http://%s:%s/api/v1/products", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	productData := dto.ProductRequest{
		Name:           "Abiba",
		Category:       "Cleaner",
		Price:          120032.23,
		AvailableStock: 999999,
	}

	productBuf, err := json.Marshal(&productData)
	s.Require().NoError(err)

	productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, productPostResp.StatusCode)
}

func (s *TestSuite) TestCreateProdcutWithoutImageId() {
	s.CleanTable()
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

	productPostUrl := fmt.Sprintf("http://%s:%s/api/v1/products", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	productData := dto.ProductRequest{
		Name:           "Abiba",
		Category:       "Cleaner",
		Price:          120032.23,
		AvailableStock: 999999,
		SupplierId:     supplierResp.Id,
	}

	productBuf, err := json.Marshal(&productData)
	s.Require().NoError(err)

	productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, productPostResp.StatusCode)
}

func (s *TestSuite) TestCreateProductWithoutSupplierId() {
	s.CleanTable()
	client := &http.Client{}
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
		ImageId:        imageResp.Id,
	}

	productBuf, err := json.Marshal(&productData)
	s.Require().NoError(err)

	productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, productPostResp.StatusCode)
}

func (s *TestSuite) TestCreateProductEmptyProduct() {
	s.CleanTable()
	productPostUrl := fmt.Sprintf("http://%s:%s/api/v1/products", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	productData := dto.ProductRequest{}

	productBuf, err := json.Marshal(&productData)
	s.Require().NoError(err)

	productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, productPostResp.StatusCode)
}

func (s *TestSuite) TestGetAllProduct() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	for _, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)
	}

	productGetUrl := fmt.Sprintf("http://%s:%s/api/v1/products", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getReq, err := http.Get(productGetUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getReq.StatusCode)
	defer getReq.Body.Close()

	getReqBuf, err := io.ReadAll(getReq.Body)
	s.Require().NoError(err)

	extractedProduct := []dto.ProductResponse{}
	err = json.Unmarshal(getReqBuf, &extractedProduct)
	s.Require().NoError(err)

	var checkArray []dto.ProductRequest

	for _, p := range extractedProduct {
		checkArray = append(checkArray, productResponseToRequest(p))
	}

	// taking into account the limit and offset
	for i := 0; i < 10; i++ {
		s.Require().Contains(checkArray, products[i])
	}
}

func (s *TestSuite) TestGetAllProductWithLimitAndOffset() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	for _, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)
	}

	productGetUrl := fmt.Sprintf("http://%s:%s/api/v1/products?limit=5&offset=2", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getReq, err := http.Get(productGetUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getReq.StatusCode)
	defer getReq.Body.Close()

	getReqBuf, err := io.ReadAll(getReq.Body)
	s.Require().NoError(err)

	extractedProduct := []dto.ProductResponse{}
	err = json.Unmarshal(getReqBuf, &extractedProduct)
	s.Require().NoError(err)

	var checkArray []dto.ProductRequest

	for _, p := range extractedProduct {
		checkArray = append(checkArray, productResponseToRequest(p))
	}

	for i := 2; i < 7; i++ {
		s.Require().Contains(checkArray, products[i])
	}
}

func (s *TestSuite) TestGetProdcutById() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	var neededId uuid.UUID
	testId := len(products) / 2

	for i, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)

		if i == testId {
			buf, err := io.ReadAll(productPostResp.Body)
			s.Require().NoError(err)

			var temp dto.ClientResponse
			err = json.Unmarshal(buf, &temp)
			s.Require().NoError(err)
			neededId = temp.Id
		}
	}

	productGetUrl := fmt.Sprintf("http://%s:%s/api/v1/products/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())
	getReq, err := http.Get(productGetUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getReq.StatusCode)
	defer getReq.Body.Close()

	getReqBuf, err := io.ReadAll(getReq.Body)
	s.Require().NoError(err)

	extractedProduct := dto.ProductResponse{}
	err = json.Unmarshal(getReqBuf, &extractedProduct)
	s.Require().NoError(err)

	check := productResponseToRequest(extractedProduct)

	s.Require().Equal(products[testId], check)
}

func (s *TestSuite) TestGetProdcutByIdInvalidID() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	for _, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)
	}

	productGetUrl := fmt.Sprintf("http://%s:%s/api/v1/products/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, "123")
	getReq, err := http.Get(productGetUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, getReq.StatusCode)
}

func (s *TestSuite) TestGetProdcutByIdNotContained() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	for _, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)
	}

	productGetUrl := fmt.Sprintf("http://%s:%s/api/v1/products/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, uuid.New())
	getReq, err := http.Get(productGetUrl)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, getReq.StatusCode)
}

func (s *TestSuite) TestUpdateStock() {
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
	products := make([]dto.ProductRequest, 15)

	for i := 0; i < 15; i++ {
		products[i] = dto.ProductRequest{
			Name:           fmt.Sprintf("Abiba %02d", i+1),
			Category:       "Cleaner",
			Price:          120032.23,
			AvailableStock: 999999,
			SupplierId:     supplierResp.Id,
			ImageId:        imageResp.Id,
		}
	}

	var neededId uuid.UUID
	testId := len(products) / 2

	for i, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)

		if i == testId {
			buf, err := io.ReadAll(productPostResp.Body)
			s.Require().NoError(err)

			var temp dto.ClientResponse
			err = json.Unmarshal(buf, &temp)
			s.Require().NoError(err)
			neededId = temp.Id
		}
	}

	for _, p := range products {
		productBuf, err := json.Marshal(&p)
		s.Require().NoError(err)

		productPostResp, err := http.Post(productPostUrl, "application/json", bytes.NewReader(productBuf))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, productPostResp.StatusCode)
	}

	productPatchUrl := fmt.Sprintf("http://%s:%s/api/v1/products/%s?decrease=10", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())

	patchReq, err := http.NewRequest(http.MethodPatch, productPatchUrl, nil)
	s.Require().NoError(err)
	patchResp, err := client.Do(patchReq)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	getProductUrl := fmt.Sprintf("http://%s:%s/api/v1/products/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())
	getProdcutResp, err := http.Get(getProductUrl)
	s.Require().NoError(err)
	defer getProdcutResp.Body.Close()

	patchBuf, err := io.ReadAll(getProdcutResp.Body)
	s.Require().NoError(err)
	var checkProduct dto.ProductResponse
	err = json.Unmarshal(patchBuf, &checkProduct)
	s.Require().NoError(err)

	s.Require().Equal(products[testId].AvailableStock-10, checkProduct.AvailableStock)
}
