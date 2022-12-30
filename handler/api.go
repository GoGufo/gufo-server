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
	"plugin"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	ver "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"

	"github.com/spf13/viper"
)

func nomoduleAnswer(w http.ResponseWriter, r *http.Request) {

	var resp sf.ErrorResponse
	resp.Success = 0
	resp.Language = "eng"
	resp.TimeStamp = int(time.Now().Unix())
	errormsg := []sf.ErrorMsg{}
	errorans := sf.ErrorMsg{
		Code:    "000001",
		Message: "No such module or function",
	}
	errormsg = append(errormsg, errorans)
	resp.Error = errormsg

	answer, err := json.Marshal(resp)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("api.go: " + err.Error())
		}
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	w.Write([]byte(answer))
}

func fileAnswer(w http.ResponseWriter, r *http.Request, filepath string, filetype string, filename string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", filetype)
	//http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader([]byte(filepath))) //ServerContent for base64 files
	http.ServeFile(w, r, filepath) //ServerFile for download files

}

func moduleAnswer(w http.ResponseWriter, r *http.Request, s map[string]interface{}, errmsg []sf.ErrorMsg, t *sf.Request) {

	//we should indicate httpcode on error case only
	//important to return language in answer
	if errmsg != nil {
		var resp sf.ErrorResponse
		resp.Success = 0
		resp.Error = errmsg
		var lg = fmt.Sprintf("%s", s["lang"])
		if lg == "" {
			lg = "eng"
		}
		resp.Language = lg
		resp.TimeStamp = int(time.Now().Unix())
		httpsstatus := s["httpcode"].(int)

		if t.UID != "" {

			if viper.GetBool("api.go UID: " + t.UID) {
				sentry.CaptureMessage("DataBase Connection Error")
			} else {
				sf.SetErrorLog("api.go UID: " + t.UID)
			}
			//write session data in answer
			session := make(map[string]interface{})
			session["uid"] = t.UID
			session["isAdmin"] = t.IsAdmin
			session["sesionexp"] = t.SessionEnd
			session["completed"] = t.Completed
			session["readonly"] = t.Readonly
			resp.Session = session
		}
		answer, err := json.Marshal(resp)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog("api.go: " + err.Error())
			}
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Server", "Gufo")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpsstatus)
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
				session["readonly"] = t.Readonly
				resp.Session = session
			}

			answer, err := json.Marshal(resp)
			if err != nil {

				if viper.GetBool("server.sentry") {
					sentry.CaptureException(err)
				} else {
					sf.SetErrorLog("api.go: " + err.Error())
				}
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
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
	sf.SetLog(userip + " " + r.URL.Path + " " + r.Method)

	switch r.Method {
	case "OPTIONS":
		ProcessOPTIONS(w, r)
	case "GET":
		ProcessREQ(w, r)
	case "POST":
		ProcessREQ(w, r)
	case "PUT":
		ProcessPUT(w, r)
	default:
		ProcessREQ(w, r)

	}

}

func ProcessOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)

}

func ProcessPUT(w http.ResponseWriter, r *http.Request) {

	t := &sf.Request{Dbversion: ver.VERSIONDB}
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	module := patharray[2]

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

func ProcessREQ(w http.ResponseWriter, r *http.Request) {

	t := &sf.Request{Dbversion: ver.VERSIONDB}
	module := ""

	if r.Method == "GET" {

		//Determinate plugin name
		path := r.URL.Path
		patharray := strings.Split(path, "/")
		module = patharray[2]

		//	t = {t.Dbversion: ver.VERSIONDB}
		//t = sf.Request{Dbversion: ver.VERSIONDB}

	} else {

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

		//Determinate plugin name
		module = t.Module
		//	t.Dbversion = ver.VERSIONDB

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
	pluginname := fmt.Sprintf("plugins.%s", module)

	if !viper.IsSet(pluginname) {
		msg := fmt.Sprintf("No Module %s", module)
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

func loadmodule(w http.ResponseWriter, r *http.Request, mod string, t *sf.Request) {
	// load module

	plug, err := plugin.Open(mod)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("api.go:Open: " + err.Error())
		}
		nomoduleAnswer(w, r)
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
		nomoduleAnswer(w, r)
		return
	}
	// symbol - Checks the function signature
	addFunc, ok := plugin.(func(*sf.Request, *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request))
	if !ok {

		if viper.GetBool("server.sentry") {
			sentry.CaptureMessage("Plugin has no function")
		} else {
			sf.SetErrorLog("api.go:151: " + "Plugin has no function")
		}
		return

	}

	addition, errmsg, m := addFunc(t, r)
	moduleAnswer(w, r, addition, errmsg, m)
}
