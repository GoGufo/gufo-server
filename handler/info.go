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
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v "github.com/gogufo/gufo-server/version"
	sf "github.com/gogufo/gufodao"

	"github.com/spf13/viper"
)

func Info(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /info " + r.Method)

	//check for session

	session := len(r.Header["X-Authorization-Token"])
	sessionarray := make(map[string]interface{})
	if session != 0 {
		upsession := make(map[string]interface{})
		sf.SetErrorLog("info.go The token is: " + r.Header["X-Authorization-Token"][0])
		token := r.Header["X-Authorization-Token"][0]
		upsession = sf.UpdateSession(token)
		if upsession["error"] == nil {
			sf.SetErrorLog("info.go  no errors")

			sessionarray["uid"] = fmt.Sprint(upsession["uid"])
			sessionarray["isAdmin"] = upsession["isadmin"].(int)
			sessionarray["sesionexp"] = upsession["session_expired"].(int)
			sessionarray["completed"] = upsession["completed"].(int)
			sessionarray["readonly"] = upsession["readonly"].(int)
		}
	}

	ans := make(map[string]interface{})
	ans["version"] = v.VERSION
	ans["registration"] = viper.GetBool("settings.registration")

	var resp sf.SuccessResponse
	resp.Success = 1
	resp.Language = "eng"
	resp.TimeStamp = int(time.Now().Unix())
	resp.Data = ans
	resp.Session = sessionarray
	answer, err := json.Marshal(resp)
	if err != nil {
		sf.SetErrorLog("api.go:40 " + err.Error())
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(answer))

}
