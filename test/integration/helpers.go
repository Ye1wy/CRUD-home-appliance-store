package integration

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"
)

func createObject(data any, url string) error {
	payload, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	return nil
}

// func createFormData(file, filename, fieldname string) (*multipart.Writer, *bytes.Buffer, error) {
// 	buf, err := os.Open(file)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("[From sub func \"createFormData\"] Open file error: %w", err)
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	w, err := writer.CreateFormFile(fieldname, filename)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("[From sub func \"createFormData\"] CreateFormFile error: %w", err)
// 	}

// 	if _, err := io.Copy(w, buf); err != nil {
// 		return nil, nil, fmt.Errorf("[From sub func \"createFormData\"] Copy error from image buf to: %w", err)
// 	}

// 	writer.Close()

// 	return writer, body, nil
// }

func extractImageData(path, filename string) (*bytes.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[From sub func \"extractImageData\"] Open file error: %w", err)
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("[From sub func \"extractImageData\"] Decode image file error: %w", err)
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, nil)
	if err != nil {
		return nil, fmt.Errorf("[From sub func \"extractImageData\"] Encode image file error: %w", err)
	}

	return buf, nil
}

func hashBytes(data []byte) string {
	outputHash := sha256.Sum256(data)
	return hex.EncodeToString(outputHash[:])
}
