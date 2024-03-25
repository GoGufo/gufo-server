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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	viper "github.com/spf13/viper"
)

func GRPCGet(misroservice string, param string, paramid string, args map[string]interface{}, token string) map[string]interface{} {

	ans := make(map[string]interface{})

	erphost := viper.GetString("server.domain")

	header := "Bearer " + token
	URL := fmt.Sprintf("%s/api/v3/%s/%s", erphost, misroservice, param)
	if paramid != "" {
		URL = fmt.Sprintf("%s/%s", URL, paramid)
	}

	if len(args) != 0 {

		var b []string
		for key, value := range args {
			str := fmt.Sprintf("%s=%s", key, value)
			b = append(b, str)
		}
		URLValues := strings.Join(b, "&")
		URL = fmt.Sprintf("%s?%s", URL, URLValues)

	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		ans["error"] = err.Error()
		ans["httpcode"] = 400
		//	return ErrorReturn(t, 400, "000005", err.Error())

	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {header},
	}

	res, err := client.Do(req)
	if err != nil {
		ans["error"] = err.Error()
		ans["httpcode"] = 400
		//return ErrorReturn(t, 400, "000005", err.Error())
	}

	var cResp Response

	if err = json.NewDecoder(res.Body).Decode(&cResp); err != nil {
		//	return ErrorReturn(t, 400, "000005", err.Error())
		ans["error"] = err.Error()
		ans["httpcode"] = 400
	}

	ans["answer"] = cResp
	ans["httpcode"] = res.StatusCode

	return ans
}
