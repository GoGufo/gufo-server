// Copyright 2024 Alexey Yanchenko <mail@yanchenko.me>
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

package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	pbv "gopkg.in/cheggaaa/pb.v1"
)

const chunkSize = 64 * 1024

type uploader struct {
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	pool        *pbv.Pool
	DoneRequest chan string
	FailRequest chan string
}

type PBRequest struct {
	*pb.Request
}

func (t *PBRequest) MSCommunication(host string, port string) (answer map[string]interface{}) {

	answer = make(map[string]interface{})

	connection := fmt.Sprintf("%s:%s", host, port)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	//	args := os.Args
	conn, err := grpc.Dial(connection, opts...)

	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}

		answer["httpcode"] = 400
		answer["code"] = "0000234"
		answer["message"] = err.Error()
		return answer

	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)

	response, err := client.Do(context.Background(), t.Request)

	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}
		answer["httpcode"] = 500
		answer["code"] = "0000236"
		answer["message"] = fmt.Sprintf("Module connection error: %s", err.Error())
		return answer

	}

	answer = sf.ToMapStringInterface(response.Data)

	//update *sf.Request
	if response.RequestBack.Token != t.Token {
		t.Token = response.RequestBack.Token
	}
	if response.RequestBack.TokenType != t.TokenType {
		t.TokenType = response.RequestBack.TokenType
	}
	if response.RequestBack.Language != t.Language {
		t.Language = response.RequestBack.Language
	}
	if response.RequestBack.UID != t.UID {
		t.UID = response.RequestBack.UID
	}
	if response.RequestBack.IsAdmin != t.IsAdmin {
		t.IsAdmin = response.RequestBack.IsAdmin
	}
	if response.RequestBack.SessionEnd != t.SessionEnd {
		t.SessionEnd = response.RequestBack.SessionEnd
	}
	if response.RequestBack.Completed != t.Completed {
		t.Completed = response.RequestBack.Completed
	}
	if response.RequestBack.Readonly != t.Readonly {
		t.Readonly = response.RequestBack.Readonly
	}

	return answer

}

func connectgrpc(w http.ResponseWriter, r *http.Request, t *pb.Request) {

	st := PBRequest{}
	st.Request = t

	port := ""
	host := ""
	pluginname := fmt.Sprintf("microservices.%s", *t.Module)
	plygintype := ""

	msmethod := viper.GetBool("server.masterservice")

	if *t.Module != "masterservice" && msmethod {

		//Check masterservice for host and port
		host = viper.GetString("microservices.masterservice.host")
		port = viper.GetString("microservices.masterservice.port")

		//Save curent data
		//	curparam := *t.Param
		//	curmethod := *t.Method

		//api/v3/auth/signup
		//api/v3/masterservice/getmicroservicebypath

		//Modify data for request masterservice
		*st.Request.MS.Param = "getmicroservicebypath"
		*st.Request.MS.Method = "GET"

		ans := st.MSCommunication(host, port)
		if ans["httpcode"] != nil {
			httpcode := 0

			httpcode, _ = strconv.Atoi(fmt.Sprintf("%v", ans["httpcode"]))

			errorAnswer(w, r, t, httpcode, fmt.Sprintf("%v", ans["code"]), fmt.Sprintf("%v", ans["message"]))
			return
		}

		host = fmt.Sprintf("%v", ans["host"])
		port = fmt.Sprintf("%v", ans["port"])
		if ans["isinternal"] != nil {
			isint, _ := strconv.ParseBool(fmt.Sprintf("%v", ans["isinternal"]))
			if isint {
				plygintype = "internal"
			}
		}

		//Put previoud data back
		//	*st.Request.Param = curparam
		//	*st.Request.Method = curmethod

	} else {

		if !viper.IsSet(pluginname) {
			msg := fmt.Sprintf("No Module %s", *t.Module)
			if viper.GetBool("server.sentry") {
				sentry.CaptureMessage(msg)
			} else {
				sf.SetErrorLog(msg)
			}
			errorAnswer(w, r, t, 401, "0000235", msg)
			return
		}

		hostpath := fmt.Sprintf("%s.host", pluginname)
		portpath := fmt.Sprintf("%s.port", pluginname)
		host = viper.GetString(hostpath)
		port = viper.GetString(portpath)
	}

	if !msmethod {
		plygintype = fmt.Sprintf("%s.type", pluginname)
	}

	if plygintype == "internal" {

		if len(r.Header["X-Sign"]) == 0 {
			errorAnswer(w, r, t, 401, "0000234", "You have no rights")
			return
		}
		//Check for X-Sign
		signheader := r.Header["X-Sign"][0]
		sgn := viper.GetString("server.sign")

		if sgn != signheader {
			errorAnswer(w, r, t, 401, "0000234", "You have no rights")
			return
		}

	}

	ans := st.MSCommunication(host, port)
	t = st.Request

	moduleAnswerv3(w, r, ans, t)

}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}
