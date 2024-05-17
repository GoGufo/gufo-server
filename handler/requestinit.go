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

package handler

import (
	"net/http"
	"strings"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func RequestInit(r *http.Request) *pb.Request {
	t := &pb.Request{}
	p := bluemonday.UGCPolicy()

	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	module := p.Sanitize(patharray[3])
	t.Module = &module
	t.Path = &path
	t.Method = &r.Method

	sgn := viper.GetString("server.sign")
	curip := sf.ReadUserIP(r)
	usagent := r.UserAgent()

	t.Sign = &sgn

	t.IP = &curip

	t.UserAgent = &usagent

	//Function in Plugin
	if pathlenth >= 5 {
		ptr := p.Sanitize(patharray[4])
		t.Param = &ptr
	}

	//ID for function in plugin
	if pathlenth >= 6 {
		ptrs := p.Sanitize(patharray[5])
		t.ParamID = &ptrs

	}

	return t
}
