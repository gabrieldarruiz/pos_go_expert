package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Cotacao struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("Servidor iniciado na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")
	select {
	case <-time.After(5 * time.Second):
		log.Println("Request processada com sucesso")
		w.Write([]byte("Request processada com sucesso"))
	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	}
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(ctxAPI, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Erro criando request externa", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao chamar API externa:", err)
		http.Error(w, "Erro na API externa", http.StatusGatewayTimeout)
		return
	}
	defer resp.Body.Close()

	var c Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		http.Error(w, "Erro no JSON", http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer cancelDB()

	err = insertCotacao(ctxDB, c.USDBRL.Bid)
	if err != nil {
		log.Println("Erro ao salvar no banco:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": c.USDBRL.Bid})
}

func insertCotacao(ctx context.Context, bid string) error {
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	return err
}
