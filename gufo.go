// Copyright 2020-2024 Alexey Yanchenko <mail@yanchenko.me>
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
// This is main file, from which starts Gufo

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/certifi/gocertifi"
	handler "github.com/gogufo/gufo-api-gateway/handler"
	v "github.com/gogufo/gufo-api-gateway/version"

	"github.com/getsentry/sentry-go"
	viper "github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

// info is function for CLI.
// in this function determinate to start Web server
func info() {

	app.Name = "Gufo API Gateway"
	app.Usage = "RESTFull API with GRPC microservices"
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

	if viper.GetBool("server.sentry") {

		sf.SetLog("Connect to Setry...")

		sentryClientOptions := sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			EnableTracing:    viper.GetBool("sentry.tracing"),
			Debug:            viper.GetBool("sentry.debug"),
			TracesSampleRate: viper.GetFloat64("sentry.trace"),
		}

		rootCAs, err := gocertifi.CACerts()
		if err != nil {
			sf.SetLog("Could not load CA Certificates for Sentry: " + err.Error())

		} else {
			sentryClientOptions.CaCerts = rootCAs
		}

		err = sentry.Init(sentryClientOptions)

		if err != nil {
			sf.SetLog("Error with sentry.Init: " + err.Error())
		}

		flushsec := viper.GetDuration("sentry.flush")

		defer sentry.Flush(flushsec * time.Second)

	}

	// run CLI function
	info()
	commands()
	flags()

	// Run CLI + Web Server
	err := app.Run(os.Args)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("gufo.go:101: " + err.Error())
		}
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

// ExitApp is Handler function for stop app by GET requet. Works in Debug mode only
func ExitApp(w http.ResponseWriter, r *http.Request) {
	var m string = "Gufo Stop \t"
	sf.SetLog(m)
	fmt.Printf(m)
	os.Exit(3)
}

// StartService is function for start WEB Server to listen port
func StartService(c *cli.Context) (rtnerr error) {
	//Initiate redis cache
	sf.InitCache()
	port := sf.ConfigString("server.port")

	var m string = "Gufo Starting. Listen []:" + port + "\t"
	sf.SetLog(m)
	fmt.Printf(m)
	http.HandleFunc("/api/", handler.WrongRequest)
	http.HandleFunc("/api/v3/", func(w http.ResponseWriter, r *http.Request) { handler.API(w, r, 3) })
	http.HandleFunc("/api/v3/info", handler.Info)
	http.HandleFunc("/api/v3/health", handler.Health)

	if viper.GetBool("server.debug") {
		http.HandleFunc("/exit", ExitApp)
	}

	//Server start
	//go http.ListenAndServe(":"+port, nil)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("gufo.go: " + err.Error())
		}

		os.Exit(1)
	}

	return nil
}
