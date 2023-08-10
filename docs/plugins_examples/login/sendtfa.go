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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func otp(t *sf.Request, r *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {
	ans := make(map[string]interface{})

	uname := t.ParamID

	var userExist sf.Users
	var userExistInfo UsersInfo

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

	rows := db.Conn.Where(`(name = ? OR mail = ?)`, uname, uname).First(&userExist)
	db.Conn.Where(`(uid = ?)`, userExist.UID).First(&userExistInfo)

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

	go sendtfa(t, userExistInfo.Name, userExist.Mail, uname)

	ans["2fa"] = true

	return ans, nil, t

}

func sendtfa(t *sf.Request, helloname string, email string, uname string) {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	sysdir := viper.GetString("server.sysdir")
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/en.toml", sysdir))
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/ru.toml", sysdir))

	otp := sf.Numgen(6)

	//Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("signup.go: " + err.Error())
		}

		return
	}

	timehash := sf.TimeHash{
		UID:      uname,
		Hash:     otp,
		Mail:     email,
		Param:    "otp",
		Created:  int(time.Now().Unix()),
		Livetime: int(time.Now().Unix()) + 300,
	}
	db.Conn.Create(&timehash)

	lg := ""
	switch t.Language {
	case "english":
		lg = "eng"
	case "russian":
		lg = "ru"
	case "deutsche":
		lg = "de"
		/*
		   case "italian":
		     lg = "it"
		   case "franch":
		     lg = "fr"
		   case "spanish":
		     lg = "sp" */
	default:
		lg = "eng"
	}

	localizer := i18n.NewLocalizer(bundle, lg)

	link := fmt.Sprintf("%s", otp)
	subj := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Subject3",
		},
	})

	title := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PassEmailTitle",
		},
		TemplateData: map[string]string{
			"Name": helloname,
		},
	})

	conflinka := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "OTP",
		},
		TemplateData: map[string]string{
			"Password": link,
		},
	})

	linkarray := []string{conflinka}

	ms := &sf.MailSettings{}
	ms.Custom = false
	go sf.SendHTMLEmail(email, title, linkarray, subj, "email.html", nil, ms)

}
