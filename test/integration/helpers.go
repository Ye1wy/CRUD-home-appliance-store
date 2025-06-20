package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
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
