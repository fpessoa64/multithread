package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fpessoa64/multithread/internal/types"
)

const BrasilApiName = "brasilapi"

type dataBrasilApi struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type BrasilApi struct {
}

func NewBrasilApi() *BrasilApi {
	return &BrasilApi{}
}

func (b *BrasilApi) ToString(result types.Result) string {

	var data dataBrasilApi
	jsonBytes, err := json.Marshal(result.Data)
	if err != nil {
		return fmt.Sprintf("Error marshaling data: %v", err)
	}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return fmt.Sprintf("Error parsing data: %v", err)
	}

	return fmt.Sprintf("%s: %v", result.Source, data)
}

func (b *BrasilApi) Fetch(cep string, ch chan<- types.Result) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ch <- types.Result{
			Source: BrasilApiName,
			Err:    err,
		}
		return
	}

	ch <- types.Result{
		Source: BrasilApiName,
		Data:   data,
	}
}
