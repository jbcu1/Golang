package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)
type Users struct {
	Users []User
}

// var info []Info
// a type and a list of social links
type User struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"Age"`
	Social Social `json:"social"`
}

// Social struct which contains a
// list of links
type Social struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

func main(){


	jsonF,err:=os.Open("/home/jbcui/go/src/prox/out.json")
	if err!=nil{
		panic(err)
	}
	defer jsonF.Close()
	byteV,_:=ioutil.ReadAll(jsonF)
	var info Users
	json.Unmarshal(byteV, &info)

	for i:=range info.Users{
		fmt.Println(info.Users[i].Type)
	}

}