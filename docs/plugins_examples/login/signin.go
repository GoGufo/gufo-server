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
	"fmt"
	"net/http"
	"strings"
	"time"

	"erp/model"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"
	"github.com/spf13/viper"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
)

// POST only
func Signin(t *sf.Request, r *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()
	errormsg := []sf.ErrorMsg{}

	ottoken := ""

	if r.URL.Query().Get("ot_token") != "" {
		ottoken = p.Sanitize(r.URL.Query().Get("ot_token"))
	}

	if ottoken != "" {

		ans, errormsg, t = signinwithtoken(t, r)
		return ans, errormsg, t
	}

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

	//1. Check for need  data
	creds := &SignInCred{}
	var userExist sf.Users

	ptfa := p.Sanitize(fmt.Sprintf("%v", t.Args["tfa"]))

	if ptfa != "" {
		// Check 2FA
		tfa := ptfa
		uname := p.Sanitize(fmt.Sprintf("%v", t.Args["user"]))
		timehash := sf.TimeHash{}

		rows := db.Conn.Debug().Where(`(uid = ? AND hash = ?)`, uname, tfa).First(&timehash)

		if rows.RowsAffected == 0 {
			// return error. user name is exist in db users
			ans["httpcode"] = 400
			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000021",
				Message: "There is no data",
			}
			errormsg = append(errormsg, errorans)
			return ans, errormsg, t
		}

		// Check for OTP livetime
		ctime := int(time.Now().Unix())

		if ctime > timehash.Livetime {
			//Delete OTP
			db.Conn.Delete(sf.TimeHash{}, "uid = ? AND hash = ?", uname, tfa)

			ans["httpcode"] = 400
			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000022",
				Message: "You already asked for confirmation email",
			}
			errormsg = append(errormsg, errorans)

			return ans, errormsg, t
		}

		//If right - return token

		//	var userExistInfo model.UsersInfo

		rows = db.Conn.Debug().Where(`(name = ? OR mail = ?)`, uname, uname).First(&userExist)

		if rows.RowsAffected == 0 {
			// return error. user name is exist in db users

			ans["httpcode"] = 400
			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000003",
				Message: "There is no such user",
			}
			errormsg = append(errormsg, errorans)

			return ans, errormsg, t
		}

		//2. If user active
		if userExist.Status == 0 {

			ans["httpcode"] = 400

			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000013",
				Message: "User blocked",
			}
			errormsg = append(errormsg, errorans)
			return ans, errormsg, t
		}

		//4. Check if user confirmed his email
		ans["email_confirmed"] = true
		if userExist.Completed == 0 {
			ans["email_confirmed"] = false //User blocked

		}

		//5. Create token
		token, expecttime, err := sf.SetSession(userExist.UID, userExist.IsAdmin, userExist.Completed, userExist.Readonly)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog("signin.go: " + err.Error())
			}
		}

		t.IsAdmin = 0
		if userExist.IsAdmin == 1 {
			t.IsAdmin = 1
		}

		//6. Write data to users table

		db.Conn.Table("users").Where("uid = ?", userExist.UID).Updates(map[string]interface{}{"access": int(time.Now().Unix()), "login": int(time.Now().Unix())})

		//7. Get special information
		var userInfo model.UsersInfo
		db.Conn.Where(`(uid = ?)`, userExist.UID).First(&userInfo)
		ans["companyid"] = userInfo.CompanyID

		//8. TODO Write data to signin history table

		//return data

		t.UID = userExist.UID
		t.SessionEnd = expecttime

		ans["token"] = token
		//ans["uid"] = userExist.UID
		ans["username"] = userExist.Name
		ans["email"] = userExist.Mail
		//ans["session_expired"] = expecttime

		return ans, nil, t

	}

	creds.Username = p.Sanitize(fmt.Sprintf("%v", t.Args["user"]))
	creds.Password = p.Sanitize(fmt.Sprintf("%v", t.Args["pass"]))

	if creds.Username == "" || creds.Password == "" {

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000001",
			Message: "Missing Name or Password",
		}
		errormsg = append(errormsg, errorans)
		ans["httpcode"] = 400
		return ans, errormsg, t
	}

	//2. Check if user exist

	var userExistInfo model.UsersInfo

	rows := db.Conn.Debug().Where(`(name = ? OR mail = ?)`, creds.Username, creds.Username).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users

		ans["httpcode"] = 400

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000003",
			Message: "There is no such user",
		}
		errormsg = append(errormsg, errorans)
		return ans, errormsg, t
	}

	//2. If user active
	if userExist.Status == 0 {

		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000013",
			Message: "User blocked",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	//3. Check password

	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Pass), []byte(creds.Password)); err != nil {
		// Password not matched
		ans["httpcode"] = 400

		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000008",
			Message: "Password not matched",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t

	}

	db.Conn.Where(`(uid = ?)`, userExist.UID).First(&userExistInfo)

	//3.1 Check for 2FA
	if userExistInfo.TFA != 0 {

		//1. Generate OTP and send email. Retrun user information about 2FA required

		go sendtfa(t, userExistInfo.Name, userExist.Mail, creds.Username)

		askedemail := maskemail(userExist.Mail)

		ans["tfa"] = true
		ans["tfatype"] = userExistInfo.TFAType
		ans["sendto"] = askedemail

		return ans, nil, t
	}

	//4. Check if user confirmed his email
	ans["email_confirmed"] = true
	if userExist.Completed == 0 {
		ans["email_confirmed"] = false //User blocked

	}

	//5. Create token
	token, expecttime, err := sf.SetSession(userExist.UID, userExist.IsAdmin, userExist.Completed, userExist.Readonly)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("signin.go: " + err.Error())
		}
	}

	t.IsAdmin = 0
	if userExist.IsAdmin == 1 {
		t.IsAdmin = 1
	}

	//6. Write data to users table

	db.Conn.Table("users").Where("uid = ?", userExist.UID).Updates(map[string]interface{}{"access": int(time.Now().Unix()), "login": int(time.Now().Unix())})

	//7. Get special information
	var userInfo model.UsersInfo
	db.Conn.Where(`(uid = ?)`, userExist.UID).First(&userInfo)
	ans["companyid"] = userInfo.CompanyID
	ans["staffid"] = userInfo.StaffID

	//8. Get information about show or not support table
	//Check if Table InvoicingSteps exist
	if !db.Conn.Migrator().HasTable(&model.InvoicingSteps{}) {

		//Create usersdate table
		db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&model.InvoicingSteps{})
	}
	invoicingsetpsdata := &model.InvoicingSteps{}
	supporttable := 0

	xrows := db.Conn.Model(&invoicingsetpsdata).Where("uid = ?", t.UID).First(&invoicingsetpsdata)

	if xrows.RowsAffected != 0 {
		//Create Data
		supporttable = invoicingsetpsdata.ShowMe
	}

	//8. TODO Write data to signin history table

	//return data

	t.UID = userExist.UID
	t.SessionEnd = expecttime

	ans["token"] = token
	//ans["uid"] = userExist.UID
	ans["username"] = userExist.Name
	ans["email"] = userExist.Mail
	ans["support_table"] = supporttable
	//ans["session_expired"] = expecttime

	return ans, nil, t

}

func maskemail(email string) string {
	mailarr := strings.Split(email, "@")
	domain := mailarr[1]
	milbody := mailarr[0]
	domainarr := strings.Split(domain, ".")
	domainbody := domainarr[0]
	mailbodymask := milbody[0:1] + "***"
	domainbodymask := domainbody[0:1] + "***" + domainbody[len(domainbody)-1:]
	maskedemail := mailbodymask + "@" + domainbodymask + "." + domainarr[1]
	return maskedemail
}
