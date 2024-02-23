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

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func errorAnswer(w http.ResponseWriter, r *http.Request, t *pb.Request, httpcode int, code string, message string) {

	ans := make(map[string]interface{})

	ans["httpcode"] = httpcode
	ans["code"] = code
	ans["message"] = message

	moduleAnswerv3(w, r, ans, t)

}
