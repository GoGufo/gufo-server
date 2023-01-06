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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	ver "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"
	"github.com/microcosm-cc/bluemonday"

	"github.com/spf13/viper"
)

func ProcessREQ(w http.ResponseWriter, r *http.Request) {

	t := &sf.Request{Dbversion: ver.VERSIONDB}
	//Determinate plugin name, params etc.
	//
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	p := bluemonday.UGCPolicy()

	if pathlenth < 3 {

		nomoduleAnswer(w, r)
		return

	}
	//Plagin Name
	t.Module = p.Sanitize(patharray[3])

	//Function in Plugin
	if pathlenth >= 5 {
		t.Param = p.Sanitize(patharray[4])
	}

	//ID for function in plugin
	if pathlenth >= 6 {
		t.ParamID = p.Sanitize(patharray[5])

	}

	if r.Method == "POST" {

		//Decode request
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&t)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog(err.Error())
			}

		}

	}

	//check for session
	session := len(r.Header["Authorization"])

	if session != 0 {
		resp := make(map[string]interface{})
		tokenheader := r.Header["Authorization"][0]
		tokenarray := strings.Split(tokenheader, " ")
		t.Token = tokenarray[1]

		resp = sf.UpdateSession(t.Token)
		if resp["error"] == nil {
			t.UID = fmt.Sprint(resp["uid"])
			t.IsAdmin = resp["isadmin"].(int)
			t.SessionEnd = resp["session_expired"].(int)
			t.Completed = resp["completed"].(int)
			t.Readonly = resp["readonly"].(int)
		} else {
			t.UID = ""
			t.IsAdmin = 0
		}
	}

	if t.UID != "" && t.Readonly == 1 {
		nomoduleAnswer(w, r)
		return
	}

	mdir := viper.GetString("server.plugindir")
	pluginname := fmt.Sprintf("plugins.%s", t.Module)

	if !viper.IsSet(pluginname) {
		msg := fmt.Sprintf("No Module %s", t.Module)
		if viper.GetBool("server.sentry") {
			sentry.CaptureMessage(msg)
		} else {
			sf.SetErrorLog(msg)
		}
		nomoduleAnswer(w, r)
		return
	}

	file := viper.GetString(fmt.Sprintf("%s.file", pluginname))
	mod := fmt.Sprintf("%s%s", mdir, file)
	loadmodule(w, r, mod, t)

}
