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
// This file content Handler for API
// Each API function is independend plugin
// and API get reguest in connect with plugin
// Get response from plugin and answer to client
// All data is in JSON format

package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/spf13/viper"
)

func ProcessREQ(w http.ResponseWriter, r *http.Request, t *pb.Request, version int) {

	//Determinate plugin name, params etc.
	//
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	//p := bluemonday.UGCPolicy()

	if pathlenth < 3 {

		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Lenth")

		return

	}
	//Plagin Name

	vrs := "v3"
	t.APIVersion = &vrs

	if version == 2 {
		vrs = "v2"
		t.APIVersion = &vrs
	}

	if *t.Module == "entrypoint" {
		errorAnswer(w, r, t, 401, "0000235", "Wrong module")
		return
	}

	if r.Method == "POST" || r.Method == "PATCH" {

		//Decode request
		decoder := json.NewDecoder(r.Body)
		args := make(map[string]interface{})
		err := decoder.Decode(&args)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog(err.Error())
			}
			errorAnswer(w, r, t, 401, "0000238", "Can not decode POST body")
			return

		}
		t.Args = sf.ToMapStringAny(args)

	}

	if r.Method == "GET" && r.URL.Query() != nil || r.Method == "TRACE" && r.URL.Query() != nil || r.Method == "HEAD" && r.URL.Query() != nil {
		paramMap := make(map[string]interface{}, 0)
		for k, v := range r.URL.Query() {
			if len(v) == 1 && len(v[0]) != 0 {
				paramMap[k] = v[0]
			}
		}
		anydt := sf.ToMapStringAny(paramMap)
		t.Args = anydt

	}

	//check for session
	if viper.GetBool("server.session") {
		t = checksession(t, r)

		if t.UID != nil && *t.Readonly == int32(1) {

			errorAnswer(w, r, t, 401, "0000235", "Read Only User")
			return

		}
	}

	//Load microservice
	if *t.Module == "info" {
		Info(w, r, t)
		return
	}
	connectgrpc(w, r, t)

}
