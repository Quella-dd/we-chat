package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Student struct {
	Name    string
	Age     int
	Address string
}

type Teacher struct {
	Name    string
	Age     int
	Address string
	Worker  string
}

func main() {
	s := `{"Name": "admin", "Age": 20, "Address": "xian", "Worker":"teacher" }`

	// var aa interface{}
	var aa Student
	if err := json.NewDecoder(strings.NewReader(s)).Decode(&aa); err != nil {
		panic(err)
	}
	fmt.Println(aa)
	// switch v := aa.(type) {
	// case Student:
	// 	fmt.Println(v, "Student")
	// case Teacher:
	// 	fmt.Println(v, "Teacher")
	// default:
	// 	fmt.Println(v)
	// }
}
