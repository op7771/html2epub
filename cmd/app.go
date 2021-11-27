package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/op7771/aozora/html2epub"
	"os"
)

func main() {
	file, err := os.Open("/Users/nhn/Downloads/list_person_all_extended_utf8.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(file))

	rows, _ := reader.ReadAll()

	for i, row := range rows {
		if i > 0 {
			id := row[0]
			title := row[1]
			author := row[15] + row[16]
			url := row[50]
			person := row[14]
			grant := row[26]
			if "なし" == grant {
				fmt.Printf("%s,%s,%s,%s,%s,%s\n", id, title, author, url, person, grant)
				html2epub.Load(url, title, author, id, person)
			}
		}
	}
}
