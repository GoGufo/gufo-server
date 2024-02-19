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
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Interfacetoresponse(request *pb.Request, answer map[string]interface{}) (response *pb.Response) {

	decanswer := ToMapStringAny(answer)
	response = &pb.Response{
		Data:        decanswer,
		RequestBack: request,
	}

	return response

}

func ErrorReturn(t *pb.Request, httpcode int, code string, message string) (response *pb.Response) {

	ans := make(map[string]interface{})

	ans["httpcode"] = httpcode
	ans["code"] = code
	ans["message"] = message
	response = Interfacetoresponse(t, ans)

	return response
}
