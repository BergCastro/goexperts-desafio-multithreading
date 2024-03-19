package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	cep := "01153000"
	timeout := time.Duration(1 * time.Second)
	ch := make(chan string)

	
	go func() {
		client := http.Client{
			Timeout: timeout,
		}
		url := "https://brasilapi.com.br/api/cep/v1/" + cep
		resp, err := client.Get(url)
		if err != nil {
			ch <- "Timeout ou erro do BrasilAPI"
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- "Erro ao ler a resposta de BrasilAPI"
			return
		}
		ch <- "BrasilAPI: " + string(body)
	}()

	go func() {
		client := http.Client{
			Timeout: timeout,
		}
		url := "http://viacep.com.br/ws/" + cep + "/json/"
		resp, err := client.Get(url)
		if err != nil {
			ch <- "Timeout ou erro do ViaCEP"
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- "Erro ao ler a resposta de  ViaCEP"
			return
		}
		ch <- "ViaCEP: " + string(body)
	}()

	select {
	case res := <-ch:
		fmt.Println("Resposta mais rápida:", res)
	case <-time.After(timeout):
		fmt.Println("Timeout! Nehuma resposta em 1 segundo.")
	}
}