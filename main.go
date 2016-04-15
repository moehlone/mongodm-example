package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moehlone/mongodm_sample/models"

	"github.com/moehlone/mongodm_sample/controllers"
	_ "github.com/moehlone/mongodm_sample/routers"
	"github.com/zebresel-com/mongodm"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
	}

	// Specify localisation file for automatic validation output
	file, err := ioutil.ReadFile("locals/locals.json")

	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal JSON to map
	var localMap map[string]map[string]string
	json.Unmarshal(file, &localMap)

	// Configure the mongodm connection and specify localisation map
	dbConfig := &mongodm.Config{
		DatabaseHost: "127.0.0.1",
		DatabaseName: "mongodm_sample",
		Locals:       localMap["en-US"],
	}

	// Connect and check for error
	db, err := mongodm.Connect(dbConfig)

	if err != nil {

		fmt.Println("Database connection error: %v", err)

	} else {

		controllers.Database = db
	}

	// See https://godoc.org/github.com/zebresel-com/mongodm#Connection.Register
	db.Register(&models.User{}, "users")
	db.Register(&models.Message{}, "messages")

	// Start the webserver
	beego.Run()
}
