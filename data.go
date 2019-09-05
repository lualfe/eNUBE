package main

import "time"

//CloudData representa os dados extraídos de um arquivo CSV ou txt
type CloudData struct {
	ID         int
	Nome       string
	Inicio     time.Time
	Fim        time.Time
	Quantidade int
	Unidade    string
	Preco      float64
}

//UploadData envia os dados do arquivo para a database
func (data *CloudData) UploadData() (err error) {
	statement := "INSERT INTO dados (nome, inicio, fim, quantidade, unidade, preco) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(data.Nome, data.Inicio, data.Fim, data.Quantidade, data.Unidade, data.Preco).Scan(&data.ID)
	return
}

//GetData puxa informações do banco de dados para montagem de listas
func GetData(start, end time.Time) (data []CloudData, err error) {
	rows, err := Db.Query("SELECT DISTINCT ON(nome) id, nome, inicio, fim, quantidade, unidade, preco FROM dados WHERE inicio >= $1 AND fim < $2", start, end)
	if err != nil {
		return
	}
	for rows.Next() {
		line := CloudData{}
		err = rows.Scan(&line.ID, &line.Nome, &line.Inicio, &line.Fim, &line.Quantidade, &line.Unidade, &line.Preco)
		data = append(data, line)
	}
	rows.Close()
	return
}

//CountValues retorna o valor total de todos os recursos utilizados
func CountValues(start, end time.Time) (val float64, err error) {
	var value float64
	err = Db.QueryRow("SELECT SUM(preco) FROM dados WHERE inicio >= $1 AND fim < $2", start, end).Scan(&value)
	if err != nil {
		return
	}
	return value, err
}

//GroupResources agrupa o valor total gasto em cada recurso
func GroupResources(start, end time.Time) (data []CloudData, err error) {
	rows, err := Db.Query("SELECT nome, SUM(preco) FROM dados WHERE inicio >= $1 AND fim < $2 GROUP BY nome", start, end)
	if err != nil {
		return
	}
	for rows.Next() {
		group := CloudData{}
		err = rows.Scan(&group.Nome, &group.Preco)
		if err != nil {
			return
		}
		data = append(data, group)
	}
	rows.Close()
	return
}
