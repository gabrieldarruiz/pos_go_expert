package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(fmt.Errorf("erro ao criar requisição: %w", err))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("erro na requisição: %w", err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("resposta com status inesperado: %s", res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(fmt.Errorf("erro ao ler resposta: %w", err))
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(fmt.Errorf("erro ao criar arquivo: %w", err))
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("Dolar: %s", string(body))))
	if err != nil {
		panic(fmt.Errorf("erro ao escrever no arquivo: %w", err))
	}

	fmt.Println("Cotação salva em cotacao.txt com sucesso!")
}
