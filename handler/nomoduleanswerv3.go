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
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/spf13/viper"
)

func nomoduleAnswerv3(w http.ResponseWriter, r *http.Request) {

	var resp sf.Response
	resp.Language = "eng"
	/*
		if t.Language != nil {
			resp.Language = t.Language
		}
	*/
	resp.TimeStamp = int(time.Now().Unix())

	errormsg := []sf.ErrorMsg{}
	errorans := sf.ErrorMsg{
		Code:    "000001",
		Message: "No such module or function",
	}
	errormsg = append(errormsg, errorans)

	data := make(map[string]interface{})

	data["error"] = errormsg

	resp.Data = data

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
