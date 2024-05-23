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

package gufodao

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/grpc"
)

func GRPCConnect(host string, port string, t *pb.Request) (answer map[string]interface{}) {

	answer = make(map[string]interface{})

	if host == "" || port == "" {
		answer["httpcode"] = 500
		answer["code"] = "0000238"
		answer["message"] = "Host or Port not determinated"
		return answer
	}

	SetErrorLog(fmt.Sprintf("%s:%s", host, port))

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
			SetErrorLog("connectgrpc: " + err.Error())
		}

		answer["httpcode"] = 400
		answer["code"] = "0000234"
		answer["message"] = err.Error()
		return answer

	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)

	response, err := client.Do(context.Background(), t)

	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("connectgrpc: " + err.Error())
		}
		answer["httpcode"] = 500
		answer["code"] = "0000236"
		answer["message"] = fmt.Sprintf("Module connection error: %s", err.Error())
		return answer

	}

	answer = ToMapStringInterface(response.Data)

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
