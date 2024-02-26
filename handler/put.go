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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/spf13/viper"
)

func ProcessPUT(w http.ResponseWriter, r *http.Request, t *pb.Request, version int) {

	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)
	//p := bluemonday.UGCPolicy()

	if pathlenth < 3 {
		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Lenth")
		return
	}

	if *t.Module == "entrypoint" {
		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Lenth")
		return
	}

	//check for session
	t = checksession(t, r)

	if t.UID != nil && *t.Readonly == int32(1) {
		errorAnswer(w, r, t, 401, "0000235", "Read Only User")
		return
	}

	vrs := "v3"
	t.APIVersion = &vrs

	if version == 2 {
		vrs = "v2"
		t.APIVersion = &vrs
	}

	args := make(map[string]interface{})
	for key, value := range r.Form {

		if len(value) == 1 && len(value) != 0 {
			args[key] = value
		}
	}

	t.Args = sf.ToMapStringAny(args)

	/*
		var (
			buf        []byte
			firstChunk bool
		)
	*/
	//PUT mean file upload, so we check for file data
	file, handler, err := r.FormFile("file")

	if err != nil || file == nil || handler.Filename == "" {
		errorAnswer(w, r, t, 400, "0000235", "Missing File")
		return
	}

	defer file.Close()

	buft := bytes.NewBuffer(nil)
	if _, err := io.Copy(buft, file); err != nil {

		errorAnswer(w, r, t, 400, "0000235", err.Error())
		return
	}

	t.Filename = &handler.Filename
	t.File = buft.Bytes()

	//start uploader
	/*
		buf = make([]byte, chunkSize)
		firstChunk = true
		for {
			n, errRead := file.Read(buf)
			if errRead != nil {
				if errRead == io.EOF {
					errRead = nil
					break
				}
				errorAnswer(w, r, t, 400, "0000235", "errored while copying from file to buf")
				return
			}

			if firstChunk {
				t.Filename = &handler.Filename
				t.File = buf[:n]
				firstChunk = false
			} else {
				t.File = buf[:n]
			}
			if err != nil {
				errorAnswer(w, r, t, 400, "0000235", "failed to send chunk via stream file")
				return
			}

		}
	*/

	pluginname := fmt.Sprintf("microservices.%s", *t.Module)

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

	//Check is it plugin or GRPC server

	//Load microservice
	connectgrpc(w, r, t)

}
