//////////////////////////////////////////////////////////////////////////////////
// Copyright 2023 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package main

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

// POST only
func signinwithtoken(t *sf.Request, r *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()

	ottoken := ""

	if r.URL.Query().Get("ot_token") != "" {
		ottoken = p.Sanitize(r.URL.Query().Get("ot_token"))
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

	//Check Token
	timehash := sf.TimeHash{}

	rows := db.Conn.Debug().Where(`(hash = ? AND param = ?)`, ottoken, "OT_Auth").First(&timehash)

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
		db.Conn.Delete(sf.TimeHash{}, "hash = ?", ottoken)

		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000022",
			Message: "Token expired",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	uid := timehash.UID

	var userExist sf.Users

	//Get info about User
	rows = db.Conn.Debug().Where(`(uid = ?)`, uid).First(&userExist)

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

	//Delete OTP
	db.Conn.Delete(sf.TimeHash{}, "hash = ?", ottoken)

	//7. Get special information
	var userInfo UsersInfo
	db.Conn.Where(`(uid = ?)`, userExist.UID).First(&userInfo)
	ans["companyid"] = userInfo.CompanyID

	t.UID = userExist.UID
	t.SessionEnd = expecttime

	ans["email_confirmed"] = true
	if userExist.Completed == 0 {
		ans["email_confirmed"] = false //User blocked

	}

	ans["token"] = token
	//ans["uid"] = userExist.UID
	ans["username"] = userExist.Name
	ans["email"] = userExist.Mail
	//ans["session_expired"] = expecttime

	return ans, nil, t

}
