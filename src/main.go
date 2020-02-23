package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Pessoa ::
type Pessoa struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
	Cpf   string `json:"cpf"`
}

var pessoas = []Pessoa{
	{
		ID:    1,
		Nome:  "Cayo Andrade",
		Idade: 21,
		Cpf:   "123.456.789-10",
	},
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/pessoas", buscarTodasAsPessoas).Methods("GET")
	router.HandleFunc("/pessoas/{id}", buscarUmaPessoa).Methods("GET")
	router.HandleFunc("/pessoas/{id}", atualizarUmaPessoa).Methods("PUT")
	router.HandleFunc("/pessoas/{id}", deletarUmaPessoa).Methods("DELETE")
	router.HandleFunc("/pessoas", cadastrarNovaPessoa).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func buscarTodasAsPessoas(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pessoas)
}

func buscarUmaPessoa(w http.ResponseWriter, r *http.Request) {
	pessoaID := mux.Vars(r)["id"]

	pessoaIDInteger, err := strconv.Atoi(pessoaID)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, pessoa := range pessoas {
		if pessoa.ID == pessoaIDInteger {
			json.NewEncoder(w).Encode(pessoa)
		}
	}
}

func cadastrarNovaPessoa(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Erro ao converter dados recebidos.")
	}

	json.Unmarshal(reqBody, &pessoa)
	pessoas = append(pessoas, pessoa)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(pessoa)
}

func atualizarUmaPessoa(w http.ResponseWriter, r *http.Request) {
	var pessoaAtt Pessoa

	pessoaID := mux.Vars(r)["id"]

	pessoaIDInteger, err := strconv.Atoi(pessoaID)
	if err != nil {
		log.Fatal(err.Error())
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Erro ao converter dados recebidos.")
	}

	json.Unmarshal(reqBody, &pessoaAtt)

	for i, pessoaArr := range pessoas {
		if pessoaArr.ID == pessoaIDInteger {
			pessoaArr.Nome = pessoaAtt.Nome
			pessoaArr.Idade = pessoaAtt.Idade
			pessoaArr.Cpf = pessoaAtt.Cpf
			pessoas = append(pessoas[:i], pessoaArr)
			json.NewEncoder(w).Encode(pessoaArr)
		}
	}
}

func deletarUmaPessoa(w http.ResponseWriter, r *http.Request) {
	pessoaID := mux.Vars(r)["id"]

	pessoaIDInteger, err := strconv.Atoi(pessoaID)
	if err != nil {
		log.Fatal(err.Error())
	}

	for i, pessoaSingle := range pessoas {
		if pessoaSingle.ID == pessoaIDInteger {
			pessoas = append(pessoas[:i], pessoas[i+1:]...)
			fmt.Fprintf(w, "Pessoa deletada com sucesso!")
		}
	}

}
