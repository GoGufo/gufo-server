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

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
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

		if msmethod {
			//Ask Masterservice for Session Host
			//	st := PBRequest{}
			//	st.Request = t

			host = viper.GetString("microservices.masterservice.host")
			port = viper.GetString("microservices.masterservice.port")

			//Modify data for request masterservice

			mst := &pb.InternalRequest{}
			param := "getsessionhost"
			gt := "GET"
			mst.Param = &param
			mst.Method = &gt

			t.IR = mst
			t.Token = &tokenheader

			ans := sf.GRPCConnect(host, port, t)
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
		mstb := &pb.InternalRequest{}
		param := "checksession"
		gt := "GET"
		mstb.Param = &param
		mstb.Method = &gt

		t.IR = mstb

		//Send Authorisation token to microservice

		ans := sf.GRPCConnect(host, port, t)
		if ans["error"] != nil {
			return t
		}

		if ans["uid"] != nil {
			uid := fmt.Sprintf("%v", ans["uid"])
			t.UID = &uid
		}
		if ans["isadmin"] != nil {
			isadminint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["isadmin"]))
			isadmin32 := int32(isadminint)
			t.IsAdmin = &isadmin32
		}
		if ans["sessionend"] != nil {
			sesint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["sessionend"]))
			sesint32 := int32(sesint)
			t.SessionEnd = &sesint32
		}
		if ans["completed"] != nil {
			comint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["completed"]))
			comint32 := int32(comint)
			t.Completed = &comint32
		}
		if ans["readonly"] != nil {
			roint, _ := strconv.Atoi(fmt.Sprintf("%v", ans["readonly"]))
			roint32 := int32(roint)
			t.Readonly = &roint32
		}
		if ans["token"] != nil {
			tkn := fmt.Sprintf("%v", ans["token"])
			t.Token = &tkn
		}
		if ans["token_type"] != nil {
			tkntp := fmt.Sprintf("%v", ans["token_type"])
			t.TokenType = &tkntp
		}

	}

	return t

}
