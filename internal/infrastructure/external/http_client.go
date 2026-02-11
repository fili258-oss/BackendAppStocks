package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient es un cliente HTTP reutilizable con configuración
type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// HTTPClientConfig configuración para el cliente HTTP
type HTTPClientConfig struct {
	BaseURL string
	Timeout time.Duration
	Headers map[string]string
}

// NewHTTPClient crea una nueva instancia de HTTPClient
func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL: config.BaseURL,
		headers: config.Headers,
	}
}

// Get realiza una petición GET
func (c *HTTPClient) Get(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
	url := c.buildURL(path, queryParams)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	// Agregar headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetJSON realiza una petición GET y parsea JSON
func (c *HTTPClient) GetJSON(ctx context.Context, path string, queryParams map[string]string, result interface{}) error {
	body, err := c.Get(ctx, path, queryParams)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error al parsear JSON: %w", err)
	}

	return nil
}

// buildURL construye la URL completa con query parameters
func (c *HTTPClient) buildURL(path string, queryParams map[string]string) string {
	url := c.baseURL + path

	if len(queryParams) > 0 {
		url += "?"
		first := true
		for key, value := range queryParams {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	return url
}
