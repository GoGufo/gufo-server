// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/spf13/viper"
)

func WrongRequest(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data

	//check for session
	ans := make(map[string]interface{})
	ans["code"] = "000056"
	ans["mesage"] = "Wrong Request"
	ans["httpcode"] = 404

	var resp sf.Response
	answer, err := json.Marshal(resp)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("info.go " + err.Error())
		}
		return
	}

	for i := 0; i < len(HeaderKeys); i++ {
		w.Header().Set(HeaderKeys[i], HeaderValues[i])
	}
	w.WriteHeader(404)
	w.Write([]byte(answer))

}
