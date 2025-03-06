package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Estrutura para a resposta da BrasilAPI
type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

// Estrutura para a resposta da ViaCEP
type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

// Estrutura unificada para resposta
type CepResult struct {
	API     string
	Address interface{}
	Err     error
}

func main() {
	cep := "22450000"

	// Criar um canal para receber os resultados
	resultChannel := make(chan CepResult)

	// Criar um contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Iniciar as goroutines para cada API
	go fetchBrasilAPI(ctx, cep, resultChannel)
	go fetchViaCEP(ctx, cep, resultChannel)

	// Aguardar o primeiro resultado
	select {
	case result := <-resultChannel:
		if result.Err != nil {
			fmt.Printf("Erro na API %s: %v\n", result.API, result.Err)
			return
		}
		printResult(result)
	case <-ctx.Done():
		fmt.Println("Timeout: As APIs demoraram mais de 1 segundo para responder")
		return
	}
}

func fetchBrasilAPI(ctx context.Context, cep string, ch chan<- CepResult) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- CepResult{API: "BrasilAPI", Err: err}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- CepResult{API: "BrasilAPI", Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- CepResult{API: "BrasilAPI", Err: err}
		return
	}

	var result BrasilAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		ch <- CepResult{API: "BrasilAPI", Err: err}
		return
	}

	ch <- CepResult{API: "BrasilAPI", Address: result}
}

func fetchViaCEP(ctx context.Context, cep string, ch chan<- CepResult) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- CepResult{API: "ViaCEP", Err: err}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- CepResult{API: "ViaCEP", Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- CepResult{API: "ViaCEP", Err: err}
		return
	}

	var result ViaCEPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		ch <- CepResult{API: "ViaCEP", Err: err}
		return
	}

	ch <- CepResult{API: "ViaCEP", Address: result}
}

func printResult(result CepResult) {
	fmt.Printf("Resultado mais rápido obtido da API: %s\n", result.API)
	fmt.Println("Dados do endereço:")

	switch addr := result.Address.(type) {
	case BrasilAPIResponse:
		fmt.Printf("CEP: %s\n", addr.Cep)
		fmt.Printf("Estado: %s\n", addr.State)
		fmt.Printf("Cidade: %s\n", addr.City)
		fmt.Printf("Bairro: %s\n", addr.Neighborhood)
		fmt.Printf("Rua: %s\n", addr.Street)
	case ViaCEPResponse:
		fmt.Printf("CEP: %s\n", addr.Cep)
		fmt.Printf("Estado: %s\n", addr.Uf)
		fmt.Printf("Cidade: %s\n", addr.Localidade)
		fmt.Printf("Bairro: %s\n", addr.Bairro)
		fmt.Printf("Rua: %s\n", addr.Logradouro)
	}
}
