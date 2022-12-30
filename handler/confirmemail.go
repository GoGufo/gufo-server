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
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"
	"github.com/spf13/viper"

	"github.com/microcosm-cc/bluemonday"
)

type ConfEmailLink struct {
	Email string
	Token string
	Lang  string
}

func Confirmemail(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /confirmemail " + r.Method)

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()

	//1. We get request with email and hash and check if it need data exist
	if r.URL.Query()["token"][0] == "" || r.URL.Query()["email"][0] == "" {

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000001",
			Message: "Missing Token or Email",
		}
		errormsg = append(errormsg, errorans)
		sendanswer(w, r, ans, errormsg)
		return
	}

	//2. Clen it from any tags
	var data ConfEmailLink
	data.Email = p.Sanitize(r.URL.Query()["email"][0])
	data.Token = p.Sanitize(r.URL.Query()["token"][0])
	if r.URL.Query()["lang"][0] == "" {
		data.Lang = "en"
	} else {
		data.Lang = p.Sanitize(r.URL.Query()["lang"][0])
	}

	//4. Check is hash live
	var userHash sf.TimeHash

	//Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("confirmemail.go: " + err.Error())
		}
		//return "error with db"
	}

	//4.1. Check if hash is exist in db users
	rows := db.Conn.Where(`hash = ? and mail = ?`, data.Token, data.Email).First(&userHash)
	if rows.RowsAffected == 0 {
		// return error. Hash is not exist in db

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000008",
			Message: "Hash is not exist in db",
		}
		errormsg = append(errormsg, errorans)
		sendanswer(w, r, ans, errormsg)
		return
	}

	curtime := int(time.Now().Unix())
	if userHash.Livetime < curtime {

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000009",
			Message: "Hash is overtime",
		}
		errormsg = append(errormsg, errorans)
		sendanswer(w, r, ans, errormsg)
		return
	}

	//5. Update users table
	db.Conn.Table("users").Where("mail = ?", data.Email).Updates(map[string]interface{}{"completed": true, "mailconfirmed": curtime})

	//6. Delete hash
	db.Conn.Delete(sf.TimeHash{}, "hash = ?", data.Token)

	ans["response"] = "100002" // email confirmed
	sendanswer(w, r, ans, nil)
	return

}

func sendanswer(w http.ResponseWriter, r *http.Request, ans map[string]interface{}, errmsg []sf.ErrorMsg) {

	if ans["error"] != nil {
		var resp sf.ErrorResponse
		resp.Success = 0
		resp.Error = errmsg
		var lg = fmt.Sprintf("%s", ans["lang"])
		if lg == "" {
			lg = "eng"
		}
		resp.Language = lg
		resp.TimeStamp = int(time.Now().Unix())
		answer, err := json.Marshal(resp)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog("confirmemail.go: " + err.Error())
			}
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Server", "Gufo")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(answer))
	} else {
		var resp sf.SuccessResponse
		resp.Success = 1
		resp.Language = "eng"
		resp.TimeStamp = int(time.Now().Unix())
		resp.Data = ans
		answer, err := json.Marshal(resp)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog("confirmemail.go: " + err.Error())
			}
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Server", "Gufo")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(answer))
	}

}
