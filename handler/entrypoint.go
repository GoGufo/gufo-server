// Copyright 2023 Alexey Yanchenko <mail@yanchenko.me>
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

package handler

import (
	"fmt"
	"plugin"

	"github.com/getsentry/sentry-go"
	v "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"
	"github.com/spf13/viper"
)

// This function will be executed when Gufo is starting
func Entrypoint() {
	fmt.Printf("Check Entrypoint")

	// Check curent Version
	isstart := false

	// Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("entrypoint.go: " + err.Error())
		}
		//return "error with db"
	}

	var curentrypoint sf.Entrypoint
	var count int64
	db.Conn.Debug().Model(&curentrypoint).Where(`version = ? `, v.VERSION).Count(&count)

	if count == 0 {
		//Run user function
		isstart = true
	} else {
		db.Conn.Where(`version = ? `, v.VERSION).First(&curentrypoint)
		if !curentrypoint.Status {
			//Run user function
			isstart = true
		}
	}

	if isstart == false {
		fmt.Printf("No Entrypoint Job")
	}

	if isstart == true {
		fmt.Printf("Start Entrypoint Job")
		mdir := viper.GetString("server.plugindir")
		pluginname := fmt.Sprintf("plugins.%s", "entrypoint")
		file := viper.GetString(fmt.Sprintf("%s.file", pluginname))
		mod := fmt.Sprintf("%s%s", mdir, file)

		plug, err := plugin.Open(mod)
		if err != nil {
			fmt.Printf("No to Entrypoint plugin")
			return
		}

		plugin, err := plug.Lookup("Init")

		if err != nil {
			fmt.Printf("No to Entrypoint Init")
			return
		}

		addFunc, ok := plugin.(func())
		if !ok {
			fmt.Printf("No to Entrypoint Function")
			return

		}

		addFunc()

		//Update entrypoint tables
		upentrypoint := sf.Entrypoint{
			ID:      sf.Hashgen(8),
			Version: v.VERSION,
			Status:  true,
		}

		if curentrypoint.ID == "" {
			db.Conn.Create(&upentrypoint)
		} else {
			db.Conn.Model(&upentrypoint).Where("entrypointid = ?", curentrypoint.ID).Updates(&upentrypoint)
		}

		fmt.Printf("Entrypoint Job Done")

	}

}
