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
	"strconv"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func checksession(t *pb.Request, r *http.Request) *pb.Request {

	session := len(r.Header["Authorization"])
	p := bluemonday.UGCPolicy()
	xtoken := ""
	tokenheader := ""

	if session == 0 {
		//Check session from token in GET header

		if r.URL.Query().Get("access_token") != "" {
			xtoken = p.Sanitize(r.URL.Query().Get("access_token"))
			session = 1

		}

	}

	if session != 0 {

		tokenheader = r.Header["Authorization"][0]

		if xtoken != "" {
			tokenheader = xtoken
		}

	}

	if tokenheader != "" {
		//1. Check masterservice status
		msmethod := viper.GetBool("server.masterservice")
		port := ""
		host := ""

		st := PBRequest{}

		if msmethod {
			//Ask Masterservice for Session Host
			//	st := PBRequest{}
			//	st.Request = t

			host = viper.GetString("microservices.masterservice.host")
			port = viper.GetString("microservices.masterservice.port")

			//Modify data for request masterservice
			*st.Request.MS.Param = "getsessionhost"
			*st.Request.MS.Method = "GET"

			ans := st.MSCommunication(host, port)
			if ans["httpcode"] != nil {

				return t
			}

			host = fmt.Sprintf("%v", ans["host"])
			port = fmt.Sprintf("%v", ans["port"])

		} else {
			//Get session host from settings
			if !viper.IsSet("microservices.session") {
				return t
			}

			host = viper.GetString("microservices.session.host")
			port = viper.GetString("microservices.session.port")
		}

		//Connect to Session microservice to get session
		*st.Request.MS.Param = "checksession"

		//Send Authorisation token to microservice
		sesargs := make(map[string]interface{})

		sesargs["token"] = tokenheader

		ans := st.MSCommunication(host, port)
		if ans["httpcode"] != nil {
			return t
		}

		if ans["uid"] != nil {
			*t.UID = fmt.Sprintf("%v", ans["uid"])
		}
		if ans["isadmin"] != nil {
			isadminint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["sessionend"]))
			*t.IsAdmin = int32(isadminint)
		}
		if ans["sessionend"] != nil {
			sesint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["sessionend"]))
			*t.SessionEnd = int32(sesint)
		}
		if ans["completed"] != nil {
			comint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["completed"]))
			*t.Completed = int32(comint)
		}
		if ans["readonly"] != nil {
			roint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["readonly"]))
			*t.Readonly = int32(roint)
		}
		if ans["token"] != nil {
			*t.Token = fmt.Sprintf("%v", ans["token"])
		}
		if ans["token_type"] != nil {
			*t.TokenType = fmt.Sprintf("%v", ans["token_type"])
		}

	}

	return t

}
