package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//AJUDOU BASTANTE: https://blog.questionable.services/article/testing-http-handlers-go/

const erroStatus = "Status code esperado %v, mas o resultado encontrado foi %v."

//BeforeEachTest
func TestMain(m *testing.M) {
	inicializarBancoDeDados()
	retCode := m.Run()
	os.Exit(retCode)
}

//Para casos de erro
func testError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAll(t *testing.T) {

	//Criar usuário
	var usuario Usuario
	usuario.Nome = "Teste"
	data, _ := json.Marshal(usuario)
	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "/usuarios", r)
	testError(err, t)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated || resp.Body == nil {
		t.Error("Usuário não pode ser inserido")
	}

	body, err := ioutil.ReadAll(resp.Body)
	testError(err, t)

	err = json.Unmarshal(body, &usuario)
	testError(err, t)

	//Obter os usuários
	req, err = http.NewRequest("GET", "/usuarios", nil)
	testError(err, t)

	resp = httptest.NewRecorder()
	handler = http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || resp.Body == nil {
		t.Error("Usuários não encontrados")
	}

}

func TestGetByID(t *testing.T) {

	//Criar usuário
	var usuario Usuario
	usuario.Nome = "Teste"
	data, _ := json.Marshal(usuario)
	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "/usuarios", r)
	testError(err, t)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated || resp.Body == nil {
		t.Error("Usuário não pode ser inserido")
	}

	body, err := ioutil.ReadAll(resp.Body)
	testError(err, t)

	err = json.Unmarshal(body, &usuario)
	testError(err, t)

	url := fmt.Sprintf("/usuarios/%v", usuario.ID)

	//Obter o usuário
	req, err = http.NewRequest("GET", url, nil)
	testError(err, t)

	resp = httptest.NewRecorder()
	handler = http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || resp.Body == nil {
		t.Error("Usuário não encontrado")
	}

}

func TestInsert(t *testing.T) {

	var usuario Usuario
	usuario.Nome = "Teste"
	json, _ := json.Marshal(usuario)
	r := bytes.NewReader(json)

	req, err := http.NewRequest("POST", "/usuarios", r)
	testError(err, t)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf(erroStatus, http.StatusCreated, resp.Code)
	}

}

func TestUpdate(t *testing.T) {

	//Criar usuário
	var usuario Usuario
	usuario.Nome = "Teste"
	data, _ := json.Marshal(usuario)
	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "/usuarios", r)
	testError(err, t)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated || resp.Body == nil {
		t.Error("Usuário não pode ser inserido")
	}

	body, err := ioutil.ReadAll(resp.Body)
	testError(err, t)

	err = json.Unmarshal(body, &usuario)
	testError(err, t)

	//Alterar usuário
	nomeAlterado := "TesteAlterado"
	usuario.Nome = nomeAlterado
	data, _ = json.Marshal(usuario)
	r = bytes.NewReader(data)

	url := fmt.Sprintf("/usuarios/%v", usuario.ID)

	req, err = http.NewRequest("PUT", url, r)
	testError(err, t)

	resp = httptest.NewRecorder()
	handler = http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || resp.Body == nil {
		t.Error("Usuário não encontrado")
	}

	body, err = ioutil.ReadAll(resp.Body)
	testError(err, t)

	err = json.Unmarshal(body, &usuario)
	testError(err, t)

	if usuario.Nome != nomeAlterado {
		t.Error("Nome não foi alterado")
	}

}

func TestDelete(t *testing.T) {

	//Criar usuário
	var usuario Usuario
	usuario.Nome = "Teste"
	data, _ := json.Marshal(usuario)
	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "/usuarios", r)
	testError(err, t)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated || resp.Body == nil {
		t.Error("Usuário não pode ser inserido")
	}

	body, err := ioutil.ReadAll(resp.Body)
	testError(err, t)

	err = json.Unmarshal(body, &usuario)
	testError(err, t)

	url := fmt.Sprintf("/usuarios/%v", usuario.ID)

	req, err = http.NewRequest("DELETE", url, nil)
	testError(err, t)

	resp = httptest.NewRecorder()
	handler = http.HandlerFunc(UsuarioHandler)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent || resp.Body == nil {
		t.Error("Usuário não foi excluído")
	}

}
