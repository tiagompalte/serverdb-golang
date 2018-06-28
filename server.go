package main

import (
	"database/sql"
	"log"
	"net/http"
)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	return result
}

func inicializarBancoDeDados() {

	db, err := sql.Open("mysql", "root:@/")

	if err != nil {
		panic(err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer db.Close()

	exec(db, "create database if not exists cursogo")
	exec(db, "use cursogo")
	exec(db, "drop table if exists usuarios")
	exec(db, `create table usuarios (
		id integer auto_increment,
		nome varchar(80),
		PRIMARY KEY(id)
		)`)

	log.Println("Tabela criada com sucesso")
}

func main() {
	inicializarBancoDeDados()
	http.HandleFunc("/usuarios/", UsuarioHandler)
	log.Println("Executando...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
