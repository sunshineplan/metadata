package metadata

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sunshineplan/cipher"
)

// Server contains metadata server address and verify header.
type Server struct {
	// metadata server address
	Addr string
	// metadata server verify header name
	Header string
	// metadata server verify header value
	Value string
}

func (s *Server) get(metadata string, client *http.Client) ([]byte, error) {
	if metadata == "" {
		return nil, errors.New("metadata is empty")
	} else if s.Addr == "" {
		return nil, errors.New("metadata server address is empty")
	} else if s.Header == "" {
		return nil, errors.New("metadata server verify header name is empty")
	} else if s.Value == "" {
		return nil, errors.New("metadata server verify header value is empty")
	}

	url := s.Addr + "/" + metadata
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s: %s", url, err)
	}
	req.Header.Add(s.Header, s.Value)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no StatusOK response from %s: %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Get queries metadata from the metadata server.
func (s *Server) Get(metadata string, data any) error {
	return s.GetWithClient(metadata, data, http.DefaultClient)
}

// Decrypt queries encrypted metadata from the metadata server.
func (s *Server) Decrypt(metadata string, data any) error {
	return s.DecryptWithClient(metadata, data, http.DefaultClient)
}

// GetWithClient queries metadata from the metadata server
// with custom http.Client.
func (s *Server) GetWithClient(metadata string, data any, client *http.Client) error {
	b, err := s.get(metadata, client)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, data)
}

// DecryptWithClient queries encrypted metadata from the metadata server
// with custom http.Client.
func (s *Server) DecryptWithClient(metadata string, data any, client *http.Client) error {
	b, err := s.get(metadata, client)
	if err != nil {
		return err
	}

	var key string
	if err = s.GetWithClient("key", &key, client); err != nil {
		return err
	}

	str, err := cipher.DecryptText(base64.StdEncoding.EncodeToString([]byte(key)), string(b))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(str), data)
}
