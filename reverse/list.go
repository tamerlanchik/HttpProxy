package reverse

import (
	"fmt"
	"proxy/db"
	"strconv"
	"strings"
)

func RunList(n int, v bool) error {
	fmt.Println("Reverse/List: n=", n)

	db, err := db.Connect()
	if err != nil {
		return err
	}
	query := `SELECT * FROM (SELECT id, method, proto_schema, dest, created, header, body FROM request ORDER BY id DESC LIMIT $1) ORDER BY id`
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
		Header string
		Body string
	}{}
	var tmpl string
		tmpl = "|%5s|%7s|%8s|%30s|%20s|\n"
		fmt.Printf(tmpl, "id", "Meth", "Sche", "Destination", "Created")
	fmt.Println(strings.Repeat("-", 76))
	for rows.Next(){
		err := rows.Scan(&data.Id, &data.Method, &data.Schema, &data.Dest,
			&data.Created, &data.Header, &data.Body)
		if err != nil {
			return err
		}

		if v{
			fmt.Println(""+strings.Repeat("-", 76)+"")
			fmt.Printf(tmpl+"Header:\n %s" +"Body: %s\n",
				strconv.Itoa(data.Id), data.Method, data.Schema, data.Dest, data.Created,
				data.Header,
				data.Body)
		} else {
			fmt.Printf(tmpl, strconv.Itoa(data.Id), data.Method, data.Schema, data.Dest, data.Created)
		}

	}
	return nil
}

func RunDelete(id string) error {
	fmt.Println("Delete: id=", id)
	return nil
}
