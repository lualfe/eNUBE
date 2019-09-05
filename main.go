package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Db representa a conexão com o banco de dados
var Db *sql.DB

//Define as variáveis para o teste
func defineVariables() {
	os.Setenv("POSTGRES_USER", "")
	os.Setenv("POSTGRES_DBNAME", "")
	os.Setenv("POSTGRES_PASSWORD", "")
}

func init() {
	defineVariables()
	var err error
	Db, err = sql.Open("postgres", "user="+os.Getenv("POSTGRES_USER")+" dbname="+os.Getenv("POSTGRES_DBNAME")+" password="+os.Getenv("POSTGRES_PASSWORD")+" sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func formUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("dados")
	if err != nil {
		http.Error(w, "não foi possível processar o arquivo", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer file.Close()

	readFile, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "não foi possível ler o arquivo", http.StatusInternalServerError)
		return
	}

	fileType, _ := mimetype.Detect(readFile)
	if fileType != "text/plain" {
		http.Error(w, "só são aceitos arquivos txt e csv", http.StatusBadRequest)
		return
	}

	if fileType == "text/plain" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		tempFile, err := ioutil.TempFile(filepath.Join(exPath, "temp-files"), "*.txt")
		if err != nil {
			fmt.Println(err)
		}
		//O arquivo será removido após o processamento. Caso seja necessário manter o arquivo, basta alterar
		//a maneira de criação do file.
		defer os.Remove(tempFile.Name())
		tempFile.Write(readFile)
		filename, _ := os.Open(tempFile.Name())
		reader := csv.NewReader(filename)
		reader.Comma = ';'
		record, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "não foi possível processar os dados do arquivo", http.StatusInternalServerError)
			return
		}
		for i, item := range record {
			if i == 0 {
				continue
			}
			layout := "2006-01-02"
			inicio, err := time.Parse(layout, item[1])
			if err != nil {
				fmt.Println(err)
				http.Error(w, "erro processando a data de inicio", http.StatusInternalServerError)
				break
			}
			fim, err := time.Parse(layout, item[2])
			if err != nil {
				fmt.Println(err)
				http.Error(w, "erro processando a data final", http.StatusInternalServerError)
				break
			}

			quantidade, err := strconv.Atoi(item[3])
			if err != nil {
				http.Error(w, "erro processando a quantidade", http.StatusInternalServerError)
				break
			}
			preco := strings.Replace(item[5], ",", ".", -1)
			precoF, err := strconv.ParseFloat(preco, 64)
			if err != nil {
				http.Error(w, "erro processando o preço", http.StatusInternalServerError)
				break
			}

			cloudData := CloudData{Nome: item[0], Inicio: inicio, Fim: fim, Quantidade: quantidade, Unidade: item[4], Preco: precoF}
			err = cloudData.UploadData()
			if err != nil {
				http.Error(w, "erro salvando dados", http.StatusInternalServerError)
				break
			}
		}
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	}
}

//Os dados executados nessa função são renderizados diretamente no endpoint /dashboard
func showData(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02"
	//Os intervalos de data devem ser recebidos do usuário no futuro
	start, err := time.Parse(layout, "2019-01-01")
	if err != nil {
		http.Error(w, "intervalo inválido", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(layout, "2019-09-01")
	if err != nil {
		http.Error(w, "intervalo inválido", http.StatusBadRequest)
		return
	}

	preco, err := CountValues(start, end)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, "\nPRECO TOTAL\n----------")
	fmt.Fprintln(w, preco)
	fmt.Fprintln(w, "----------")

	group, err := GroupResources(start, end)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, "\nPRECO TOTAL POR RECURSO\n----------")
	var list []string
	for _, val := range group {
		fmt.Fprintln(w, val.Nome, val.Preco)
		list = append(list, val.Nome)
	}
	fmt.Fprintln(w, "----------")

	fmt.Fprintln(w, "\nLISTA DE RECURSOS UTILIZADOS\n----------")
	for _, val := range list {
		fmt.Fprintln(w, val)
	}
	fmt.Fprintln(w, "----------")
}

func newRouter() *mux.Router {
	mux := mux.NewRouter()
	return mux
}

func main() {
	mux := newRouter()
	mux.HandleFunc("/process", formUpload).Methods("POST")
	mux.HandleFunc("/dashboard", showData).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
