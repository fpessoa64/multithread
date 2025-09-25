package main

import (
	"flag"
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

func main() {
	var cep string
	flag.StringVar(&cep, "cep", "01153000", "CEP para consulta")
	flag.Parse()

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
