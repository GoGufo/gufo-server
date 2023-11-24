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

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
)

func APIv3(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " " + r.URL.Path + " " + r.Method)

	switch r.Method {
	case "OPTIONS":
		ProcessOPTIONS(w, r)
	case "GET":
		ProcessREQv3(w, r)
	case "POST":
		ProcessREQv3(w, r)
	case "DELETE":
		ProcessREQv3(w, r)
	case "PUT":
		ProcessPUTv3(w, r)
	default:
		ProcessOPTIONS(w, r)

	}

}
