package main

import (
	"encoding/json"
	"fmt"
	"os"
)
func parse_config()(*config){
	config_var := new(config)
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config_var)
	if err != nil{
		fmt.Println(err)
	}
	return config_var
}