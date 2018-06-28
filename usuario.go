package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//Usuario struct
type Usuario struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func connectDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:@/cursogo")
	checkErr(err)
	return db
}

//UsuarioHandler analisa o request e delega para função adequada
func UsuarioHandler(w http.ResponseWriter, r *http.Request) {

	sid := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		getByID(w, r, id)
	case r.Method == "GET":
		getAll(w, r)
	case r.Method == "POST":
		insert(w, r)
	case r.Method == "PUT":
		update(w, r, id)
	case r.Method == "DELETE":
		delete(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa :(")
	}

}

func getByID(w http.ResponseWriter, r *http.Request, id int) {

	db := connectDatabase()

	defer db.Close()

	var u Usuario
	db.QueryRow("select id, nome from usuarios where id = ?", id).Scan(&u.ID, &u.Nome)

	json, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(json))

}

func getAll(w http.ResponseWriter, r *http.Request) {

	db := connectDatabase()

	defer db.Close()

	rows, _ := db.Query("SELECT id, nome FROM usuarios")
	defer rows.Close()

	var listaUsuarios []Usuario
	for rows.Next() {
		var usuario Usuario
		rows.Scan(&usuario.ID, &usuario.Nome)
		listaUsuarios = append(listaUsuarios, usuario)
	}

	json, _ := json.Marshal(listaUsuarios)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(json))

}

func insert(w http.ResponseWriter, r *http.Request) {

	db := connectDatabase()

	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	var usuario Usuario
	err = json.Unmarshal(body, &usuario)
	checkErr(err)

	stmt, _ := db.Prepare("INSERT INTO usuarios(nome) VALUES(?)")
	res, err := stmt.Exec(usuario.Nome)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)

	db.QueryRow("select id, nome from usuarios where id = ?", id).Scan(&usuario.ID, &usuario.Nome)
	json, _ := json.Marshal(usuario)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(json))

}

func update(w http.ResponseWriter, r *http.Request, id int) {

	db := connectDatabase()

	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	var usuario Usuario
	err = json.Unmarshal(body, &usuario)
	checkErr(err)

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	stmt, _ := db.Prepare("UPDATE usuarios SET nome = ? WHERE id = ?")
	res, err := stmt.Exec(usuario.Nome, id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)

	var usr Usuario
	db.QueryRow("select id, nome from usuarios where id = ?", id).Scan(&usr.ID, &usr.Nome)
	json, _ := json.Marshal(usr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(json))

}

func delete(w http.ResponseWriter, r *http.Request, id int) {

	db := connectDatabase()

	defer db.Close()

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	stmt, _ := db.Prepare("DELETE FROM usuarios WHERE id = ?")
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
