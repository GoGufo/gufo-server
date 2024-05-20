// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package handler

import (
	"net/http"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	v "github.com/gogufo/gufo-api-gateway/version"

	"github.com/spf13/viper"
)

func Info(w http.ResponseWriter, r *http.Request, t *pb.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /info " + r.Method)

	ans := make(map[string]interface{})
	ans["version"] = v.VERSION
	ans["registration"] = viper.GetBool("settings.registration")

	moduleAnswerv3(w, r, ans, t)

}
