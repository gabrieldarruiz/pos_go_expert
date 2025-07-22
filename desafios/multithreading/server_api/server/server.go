package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	canal := make(chan string)
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanf("%s", &cep)
	fmt.Println("Buscando informações para o CEP:", cep)

	go func() {
		fmt.Println("[DEBUG] Iniciando chamada para Brasil API...")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
		if err != nil {
			canal <- "Erro ao criar requisição Brasil API: " + err.Error()
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro na Brasil API:", err.Error())
			canal <- "Erro ao chamar a Brasil API: " + err.Error()
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			canal <- "Erro ao ler resposta da Brasil API: " + err.Error()
			return
		}

		fmt.Println("Brasil API respondeu!")
		//time.Sleep(2 * time.Second) // para testar o timeout
		canal <- "Brasil API respondeu primeiro: " + string(body)
	}()

	go func() {
		fmt.Println("Iniciando chamada para ViaCEP API...")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
		if err != nil {
			canal <- "Erro ao criar requisição ViaCEP API: " + err.Error()
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro na ViaCEP API:", err.Error())
			canal <- "Erro ao chamar a ViaCEP API: " + err.Error()
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			canal <- "Erro ao ler resposta da ViaCEP API: " + err.Error()
			return
		}

		fmt.Println("ViaCEP API respondeu!")
		//time.Sleep(2 * time.Second) // para testar o timeout
		canal <- "ViaCEP API respondeu primeiro: " + string(body)
	}()

	fmt.Println("Aguardando resposta das APIs...")
	select {
	case mensagem := <-canal:
		fmt.Println("\n" + mensagem)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: nenhuma API respondeu em 1 segundo")
	}

	fmt.Println("Programa finalizado com sucesso!")
}
