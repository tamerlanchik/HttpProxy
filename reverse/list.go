package reverse

import (
	"fmt"
	"proxy/db"
	"strconv"
	"strings"
)

func RunList(n int) error {
	fmt.Println("Reverse/List: n=", n)

	db, err := db.Connect()
	if err != nil {
		return err
	}
	query := `SELECT id, method, proto_schema, dest, created FROM request ORDER BY id DESC LIMIT $1`
	rows, err := db.Query(query, n)
	if err != nil {
		return err
	}
	data := struct {
		Id int
		Method string
		Schema string
		Dest string
		Created string
	}{}
	tmpl := "|%5s|%7s|%8s|%30s|%20s|\n"
	fmt.Printf(tmpl, "id", "Meth", "Sche", "Destination", "Created")
	fmt.Println(strings.Repeat("-", 70 + 6))
	for rows.Next(){
		err := rows.Scan(&data.Id, &data.Method, &data.Schema, &data.Dest, &data.Created)
		if err != nil {
			return err
		}

		fmt.Printf(tmpl, strconv.Itoa(data.Id), data.Method, data.Schema, data.Dest, data.Created)
	}
	return nil
}

func RunDelete(id string) error {
	fmt.Println("Delete: id=", id)
	return nil
}
