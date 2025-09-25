package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fpessoa64/multithread/internal/types"
	"github.com/fpessoa64/multithread/internal/workers"
)

type Result struct {
	Source string
	Data   map[string]interface{}
	Err    error
}

func fetchBrasilAPI(cep string, ch chan<- Result) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- Result{Source: "BrasilAPI", Err: err}
		return
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- Result{Source: "BrasilAPI", Err: err}
		return
	}
	ch <- Result{Source: "BrasilAPI", Data: data}
}

func fetchViaCEP(cep string, ch chan<- Result) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- Result{Source: "ViaCEP", Err: err}
		return
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- Result{Source: "ViaCEP", Err: err}
		return
	}
	ch <- Result{Source: "ViaCEP", Data: data}
}

func main() {
	cep := "01153000"
	ch := make(chan types.Result, 2)
	var apis []types.ConsultCep
	apis = append(apis, workers.NewBrasilApi())
	apis = append(apis, workers.NewViaCepApi())
	for _, api := range apis {
		go api.Fetch(cep, ch)
	}

	timeout := time.After(1 * time.Second)

	select {
	case res := <-ch:
		if res.Err != nil {
			fmt.Printf("Erro na API %s: %v\n", res.Source, res.Err)
			return
		}
		fmt.Printf("Resposta da API %s:\n", res.Source)
		for k, v := range res.Data {
			fmt.Printf("%s: %v\n", k, v)
		}
	case <-timeout:
		fmt.Println("Erro: Timeout de 1 segundo atingido.")
	}
}
