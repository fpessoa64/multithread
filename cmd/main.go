package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fpessoa64/multithread/internal/types"
	"github.com/fpessoa64/multithread/internal/workers"
)

type Result struct {
	Source string
	Data   map[string]interface{}
	Err    error
}

// func fetchBrasilAPI(cep string, ch chan<- Result) {
// 	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		ch <- Result{Source: "BrasilAPI", Err: err}
// 		return
// 	}
// 	defer resp.Body.Close()
// 	var data map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		ch <- Result{Source: "BrasilAPI", Err: err}
// 		return
// 	}
// 	ch <- Result{Source: "BrasilAPI", Data: data}
// }

// func fetchViaCEP(cep string, ch chan<- Result) {
// 	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		ch <- Result{Source: "ViaCEP", Err: err}
// 		return
// 	}
// 	defer resp.Body.Close()
// 	var data map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		ch <- Result{Source: "ViaCEP", Err: err}
// 		return
// 	}
// 	ch <- Result{Source: "ViaCEP", Data: data}
// }

func main() {
	cep := "01153000"
	seconds := 1
	log.Printf("Iniciando consulta de CEPs... %s", cep)

	ch := make(chan types.Result, 2)
	var apis []types.ConsultCep
	apis = append(apis, workers.NewBrasilApi())
	apis = append(apis, workers.NewViaCepApi())
	for _, api := range apis {
		go api.Fetch(cep, ch)
	}
	timeout := time.After(time.Duration(seconds) * time.Second)
	log.Printf("Aguardando respostas por até %d segundo...", seconds)

	select {
	case res := <-ch:
		if res.Err != nil {
			log.Printf("Erro na API %s: %v\n", res.Source, res.Err)
			return
		}
		// Monta uma string única com todos os campos
		var campos []string
		for k, v := range res.Data {
			campos = append(campos, fmt.Sprintf("%s: %v", k, v))
		}
		log.Printf("Consulta retornada com sucesso pela API: %s:  %s", res.Source, fmt.Sprint(campos))
	case <-timeout:
		log.Printf("Erro: Timeout de %d segundo atingido.", seconds)
	}
	log.Println("Finalizando consulta de CEPs...")
}
