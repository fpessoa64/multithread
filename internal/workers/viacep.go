package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fpessoa64/multithread/internal/types"
)

const ViaCepApiName = "viacep"

type dataViaCepApi struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepApi struct {
}

func NewViaCepApi() *ViaCepApi {
	return &ViaCepApi{}
}

func (v *ViaCepApi) ToString(result types.Result) string {
	var data dataViaCepApi
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

func (v *ViaCepApi) Fetch(cep string, ch chan<- types.Result) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- types.Result{
			Source: ViaCepApiName,
			Err:    err,
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- types.Result{
			Source: ViaCepApiName,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- types.Result{
			Source: ViaCepApiName,
			Err:    err,
		}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ch <- types.Result{
			Source: ViaCepApiName,
			Err:    err,
		}
		return
	}

	ch <- types.Result{
		Source: ViaCepApiName,
		Data:   data,
	}
	return
}
