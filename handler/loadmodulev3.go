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
// This file content Handler for API
// Each API function is independend plugin
// and API get reguest in connect with plugin
// Get response from plugin and answer to client
// All data is in JSON format

package handler

import (
	"net/http"
	"plugin"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"

	"github.com/spf13/viper"
)

func loadmodulev3(w http.ResponseWriter, r *http.Request, mod string, t *sf.Request) {
	// load module

	plug, err := plugin.Open(mod)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("api.go:Open: " + err.Error())
		}
		nomoduleAnswerv3(w, r)
		return
	}

	//plugin, err := plug.Lookup(strings.Title(t.Param))
	plugin, err := plug.Lookup("Init")

	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("api.gp:Lookup: " + err.Error())
		}
		nomoduleAnswerv3(w, r)
		return
	}

	ans := make(map[string]interface{})

	addFunc, ok := plugin.(func(map[string]interface{}, *sf.Request, *http.Request) (map[string]interface{}, *sf.Request))
	if !ok {

		if viper.GetBool("server.sentry") {
			sentry.CaptureMessage("Plugin has no function")
		} else {
			sf.SetErrorLog("api.go: " + "Plugin has no function")
		}
		return

	}

	ans, m := addFunc(ans, t, r)
	moduleAnswerv3(w, r, ans, m)
}
