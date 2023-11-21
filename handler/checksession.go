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
	"fmt"
	"net/http"

	sf "github.com/gogufo/gufo-server/gufodao"
	"github.com/microcosm-cc/bluemonday"
)

func checksession(t *sf.Request, r *http.Request) *sf.Request {

	session := len(r.Header["Authorization"])
	p := bluemonday.UGCPolicy()
	xtoken := ""

	if session == 0 {
		//Check session from token in GET header

		if r.URL.Query().Get("access_token") != "" {
			xtoken = p.Sanitize(r.URL.Query().Get("access_token"))
			session = 1

		}

	}

	//Basic Authorisation
	if session != 0 {
		resp := make(map[string]interface{})
		tokenheader := r.Header["Authorization"][0]

		if xtoken != "" {
			tokenheader = xtoken
		}

		resp = sf.UpdateSession(tokenheader)
		if resp["error"] == nil {
			t.UID = fmt.Sprint(resp["uid"])
			t.IsAdmin = resp["isadmin"].(int)
			t.SessionEnd = resp["session_expired"].(int)
			t.Completed = resp["completed"].(int)
			t.Readonly = resp["readonly"].(int)
			t.Token = resp["token"].(string)
			t.TokenType = resp["token_type"].(string)
		} else {
			t.UID = ""
			t.IsAdmin = 0
		}
	}

	return t
}
