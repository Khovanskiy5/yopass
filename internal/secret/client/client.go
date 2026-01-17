package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
)

var HTTPClient = http.DefaultClient

type ServerError struct {
	err error
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("yopass server error: %s", e.err)
}

func (e *ServerError) Unwrap() error {
	return e.err
}

type serverResponse struct {
	Message string `json:"message"`
}

func Fetch(serverURL string, id string) (string, error) {
	serverURL = strings.TrimSuffix(serverURL, "/")

	resp, err := HTTPClient.Get(serverURL + "/secret/" + id)
	if err != nil {
		return "", &ServerError{err: err}
	}
	return handleServerResponse(resp)
}

func Store(serverURL string, s domain.Secret) (string, error) {
	serverURL = strings.TrimSuffix(serverURL, "/")

	var j bytes.Buffer
	if err := json.NewEncoder(&j).Encode(&s); err != nil {
		return "", fmt.Errorf("could not encode request: %w", err)
	}
	resp, err := HTTPClient.Post(serverURL+"/secret", "application/json", &j)
	if err != nil {
		return "", &ServerError{err: err}
	}
	return handleServerResponse(resp)
}

func handleServerResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	var r serverResponse
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(msg, &r); err == nil {
			msg = []byte(r.Message)
		}
		err := fmt.Errorf("unexpected response %s: %s", resp.Status, string(msg))
		return "", &ServerError{err: err}
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("could not decode server response: %w", err)
	}

	return r.Message, nil
}
