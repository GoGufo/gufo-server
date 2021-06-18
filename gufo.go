// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// This is main file, from which starts build Gufo

package main

import (
	"fmt"
	"net/http"
	"os"

	handler "github.com/gogufo/gufo-server/handler"
	v "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"

	viper "github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

//info is function for CLI.
// in this function determinate to start Web server
func info() {

	app.Name = "Gufo Framework"
	app.Usage = "API oriented General Universal Framework for create any Back End"
	app.Version = v.VERSION
	app.Action = StartService

}

// commands is function for CLI
func commands() {
	app.Commands = []*cli.Command{
		{
			Name:   "stop",
			Usage:  "Stop Gufo Server",
			Action: StopApp,
		},
	}
}

func flags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
		},
	}
}

// main is nain function in Gufo from which starts app.
// In this function we check all nessesary settings such as
// cinfig file, DB Connection, DB Structure, run CLI and start listen port for API requests
func main() {

	sf.CheckForFlags()

	sf.CheckConfig() // Check config file
	sf.SetLog("Check Database connection...")

	if sf.DBCheck() { // Check DB connection
		//DB connection ok
		sf.SetLog("Database connection... OK")

	} else {
		//DB Connection filed
		sf.SetErrorLog("DataBase Connection Error")
		sf.SetLog("Server Stop")
		fmt.Printf("DataBase Connection Error \t")
		fmt.Printf("Server Stop \t")
		os.Exit(3)
	}
	// Check System DB Structure. If it wrong or  missing - restore or create it
	// As well as create admin credentials, if user table is missing
	sf.CheckDBStructure()

	// run CLI function
	info()
	commands()
	flags()

	// Run CLI + Web Server
	err := app.Run(os.Args)
	if err != nil {
		sf.SetErrorLog("gufo.go:101: " + err.Error())
	}

}

// StopApp allows to stop Gufo by CLI command "stop"
func StopApp(c *cli.Context) (rtnerr error) {
	var m string = "Gufo Stop \t"
	sf.SetLog(m)
	fmt.Printf(m)
	os.Exit(3)
	return nil
}

//ExitApp is Handler function for stop app by GET requet. Works in Debug mode only
func ExitApp(w http.ResponseWriter, r *http.Request) {
	var m string = "Gufo Stop \t"
	sf.SetLog(m)
	fmt.Printf(m)
	os.Exit(3)
}

//StartService is function for start WEB Server to listen port
func StartService(c *cli.Context) (rtnerr error) {
	//Initiate redis cache
	sf.InitCache()
	port := sf.ConfigString("server.port")

	var m string = "Gufo Starting. Listen []:" + port + "\t"
	sf.SetLog(m)
	fmt.Printf(m)

	//API Handler
	http.HandleFunc("/api/confirmemail", handler.Confirmemail) //GET
	http.HandleFunc("/api/info", handler.Info)                 //GET
	http.HandleFunc("/api/logout", handler.Logout)             //GET
	http.HandleFunc("/api/", handler.API)

	if viper.GetBool("server.debug") {
		http.HandleFunc("/exit", ExitApp)
	}

	//Server start
	//go http.ListenAndServe(":"+port, nil)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		sf.SetErrorLog("gufo.go:148: " + err.Error())
		os.Exit(1)
	}

	return nil
}

/*
func conffunc(c *cli.Context) {
	if c != nil {
		sf.Configpath = *c
	}
	m := fmt.Sprintf("The conf is: %s and you input %s \t\n", sf.Configpath, c)
	fmt.Printf(m)

}
*/
