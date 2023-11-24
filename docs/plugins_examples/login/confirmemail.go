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
// Sign In
// SignIn function authorisate user in Gufo.
//

package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/spf13/viper"
)

func confirmemail(t *sf.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

	ans := make(map[string]interface{})

	//Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog(err.Error())
		}
		ans["httpcode"] = 500
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000027",
			Message: err.Error(),
		}
		errormsg = append(errormsg, errorans)
		return ans, errormsg, t
	}

	var userExist sf.Users
	var userdataExist UsersInfo

	rows := db.Conn.Where(`uid = ?`, t.UID).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users

		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "0000031",
			Message: "There is no such user",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	rows = db.Conn.Debug().Where(`uid = ?`, t.UID).First(&userdataExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "0000032",
			Message: "There is no such user",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	//Check if user already request for confirmatin email

	//check for hash lifetime
	ctime := int(time.Now().Unix())
	waittime := 300
	realtime := ctime - userExist.Mailsent
	//sf.SetErrorLog("ctime: " + fmt.Sprintf("%d", ctime))
	//sf.SetErrorLog("realtime: " + fmt.Sprintf("%d", userExist.Mailsent))
	if realtime < waittime {

		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000012",
			Message: "You already asked for confirmation email",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	db.Conn.Table("users").Where("uid = ?", t.UID).Updates(map[string]interface{}{"mailsent": int(time.Now().Unix()), "completed": 0})

	SendConfEmail(userExist.Mail, userExist.UID, userdataExist.Name, t.Language)

	ans["response"] = "100201"
	ans["message"] = "Confirmation message was sent to you email"

	return ans, nil, t
}
