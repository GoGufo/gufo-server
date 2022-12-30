// Copyright 2022 Alexey Yanchenko <mail@yanchenko.me>
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
	"fmt"
	"net/http"
	"strings"

	ver "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"

	"github.com/spf13/viper"
)

func APIv2(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " " + r.URL.Path + " " + r.Method)

	switch r.Method {
	case "OPTIONS":
		ProcessOPTIONS(w, r)
	case "GET":
		ProcessREQv2(w, r)
	case "POST":
		ProcessREQv2(w, r)
	case "PUT":
		ProcessPUT(w, r)
	default:
		ProcessREQv2(w, r)

	}

}

func ProcessREQv2(w http.ResponseWriter, r *http.Request) {
	sf.SetErrorLog("ProcessREQv2")
	sf.SetErrorLog("api.go:167 " + ver.VERSIONDB)
	t := &sf.Request{Dbversion: ver.VERSIONDB}
	module := ""

	//Determinate plugin name
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	module = patharray[2]

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
	pluginname := fmt.Sprintf("plugins.%s", module)

	if !viper.IsSet(pluginname) {
		msg := fmt.Sprintf("No Module %s", module)
		sf.SetErrorLog(msg)
		nomoduleAnswer(w, r)
		return
	}

	file := viper.GetString(fmt.Sprintf("%s.file", pluginname))
	mod := fmt.Sprintf("%s%s", mdir, file)
	loadmodule(w, r, mod, t)

}
