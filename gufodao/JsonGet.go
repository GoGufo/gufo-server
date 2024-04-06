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
	"fmt"
	"io"
	"net/http"
	"strings"
)

func JsonGet(url string, args map[string]interface{}, token string) ([]byte, error) {

	if len(args) != 0 {

		var b []string
		for key, value := range args {
			str := fmt.Sprintf("%s=%s", key, value)
			b = append(b, str)
		}
		URLValues := strings.Join(b, "&")
		url = fmt.Sprintf("%s?%s", url, URLValues)

	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if token != "" {
		header := "Bearer " + token
		req.Header.Add("Authorization", header)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	byteresponse, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return byteresponse, nil

}
