package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fpessoa64/multithread/internal/types"
)

const BrasilApiName = "brasilapi"

type BrasilApi struct {
}

func NewBrasilApi() *BrasilApi {
	return &BrasilApi{}
}

func (b *BrasilApi) Fetch(cep string, ch chan<- types.Result) error {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return err
	}

	ch <- types.Result{
		Source: BrasilApiName,
		Data:   data,
	}
	return nil
}
