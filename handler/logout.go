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
	"net/http"

	sf "github.com/gogufo/gufodao"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /logout " + r.Method)

	session := len(r.Header["X-Authorization-Token"])

	if session != 0 {
		//session exist
		sf.DelSession(r.Header["X-Authorization-Token"][0])
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Authorization-Token, Content-Type")
	w.Header().Set("Server", "Gufo")
	w.WriteHeader(200)
}
