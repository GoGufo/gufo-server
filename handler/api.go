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
	"bytes"
	"encoding/json"
	"fmt"
	sf "gufo/functions"
	ver "gufo/version"
	"net/http"
	"plugin"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func nomoduleAnswer(w http.ResponseWriter, r *http.Request) {

	var resp sf.ErrorResponse
	resp.Success = 0
	resp.Language = "eng"
	resp.TimeStamp = int(time.Now().Unix())
	resp.Error = "No such module or function"

	answer, err := json.Marshal(resp)
	if err != nil {
		sf.SetErrorLog("api.go:48 " + err.Error())
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(answer))
}

func fileAnswer(w http.ResponseWriter, r *http.Request, filepath string, filetype string, filename string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", filetype)
	http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader([]byte(filepath)))

}

func moduleAnswer(w http.ResponseWriter, r *http.Request, s map[string]interface{}, t *sf.Request) {

	if s["error"] != nil {
		var resp sf.ErrorResponse
		resp.Success = 0
		resp.Error = fmt.Sprintf("%s", s["error"])
		resp.Language = "eng"
		resp.TimeStamp = int(time.Now().Unix())

		if t.UID != "" {
			sf.SetErrorLog("api.go:67 UID: " + t.UID)
			//write session data in answer
			session := make(map[string]interface{})
			session["uid"] = t.UID
			session["isAdmin"] = t.IsAdmin
			session["sesionexp"] = t.SessionEnd
			session["completed"] = t.Completed
			resp.Session = session
		}
		answer, err := json.Marshal(resp)
		if err != nil {
			sf.SetErrorLog("api.go:77 " + err.Error())
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
		w.Header().Set("Server", "Gufo")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(answer))
	} else {
		if s["file"] != nil {
			var filename = s["file"].(string)
			fileAnswer(w, r, filename, s["filetype"].(string), s["filename"].(string))
		} else {
			var resp sf.SuccessResponse

			resp.Success = 1
			resp.Language = "eng"
			resp.TimeStamp = int(time.Now().Unix())
			resp.Data = s

			if t.UID != "" {
				//write session data in answer
				session := make(map[string]interface{})
				session["uid"] = t.UID
				session["isAdmin"] = t.IsAdmin
				session["Sesionexp"] = t.SessionEnd
				session["completed"] = t.Completed
				resp.Session = session
			}

			answer, err := json.Marshal(resp)
			if err != nil {
				sf.SetErrorLog("api.go:101 " + err.Error())
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
			w.Header().Set("Server", "Gufo")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(answer))
		}
	}
}

func API(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /api " + r.Method)

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var t *sf.Request
		err := decoder.Decode(&t)
		if err != nil {
			sf.SetErrorLog("api.go:141 " + err.Error())
		}
		module := t.Module

		t.Dbversion = ver.VERSIONDB

		//check for session

		session := len(r.Header["X-Authorization-Token"])

		if session != 0 {
			resp := make(map[string]interface{})
			t.Token = r.Header["X-Authorization-Token"][0]
			resp = sf.UpdateSession(t.Token)
			if resp["error"] == nil {
				t.UID = fmt.Sprint(resp["uid"])
				t.IsAdmin = resp["isadmin"].(int)
				t.SessionEnd = resp["session_expired"].(int)
				t.Completed = resp["completed"].(int)
			} else {
				t.UID = ""
				t.IsAdmin = 0
			}
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
	} else {
		nomoduleAnswer(w, r)
		return
	}
}

func loadmodule(w http.ResponseWriter, r *http.Request, mod string, t *sf.Request) {
	// load module
sf.SetErrorLog("the mod is " + mod)
	plug, err := plugin.Open(mod)
	if err != nil {
		sf.SetErrorLog("api.go:Open: " + err.Error())
		nomoduleAnswer(w, r)
		return
	}

	plugin, err := plug.Lookup(strings.Title(t.Param))

	if err != nil {
		sf.SetErrorLog("api.gp:Lookup: " + err.Error())
		nomoduleAnswer(w, r)
		return
	}
	// symbol - Checks the function signature
	addFunc, ok := plugin.(func(*sf.Request) (map[string]interface{}, *sf.Request))
	if !ok {
		sf.SetErrorLog("api.go:151: " + "Plugin has no function")
		return

	}

	addition, m := addFunc(t)
	moduleAnswer(w, r, addition, m)
}
