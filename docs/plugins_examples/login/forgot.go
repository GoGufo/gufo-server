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
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"

	"github.com/BurntSushi/toml"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"
)

func Forgot(t *sf.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

	ans := make(map[string]interface{})

	//1. Check for need  data

	p := bluemonday.UGCPolicy()
	email := p.Sanitize(fmt.Sprintf("%v", t.Args["email"]))
	lang := p.Sanitize(t.Language)

	if email == "" {
		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000001",
			Message: "Missing email",
		}
		errormsg = append(errormsg, errorans)
		return ans, errormsg, t
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	sysdir := viper.GetString("server.sysdir")
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/en.toml", sysdir))
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/ru.toml", sysdir))

	// send email
	lg := ""
	switch lang {
	case "english":
		lg = "eng"
	case "russian":
		lg = "ru"
		/* case "germain":
			lg = "de"
		case "italian":
			lg = "it"
		case "franch":
			lg = "fr"
		case "spanish":
			lg = "sp" */
	default:
		lg = "eng"
	}

	//2. Check if user exist
	var userExist sf.Users

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

	rows := db.Conn.Where(`mail = ?`, email).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000003",
			Message: "User is not exist",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	hashedkey := p.Sanitize(fmt.Sprintf("%v", t.Args["key"]))

	if hashedkey == "" {
		//If no confirmation code - just send this code to email
		hashkey := sf.Numgen(6)

		//Write key to key Table

		timehash := sf.TimeHash{
			UID:      userExist.UID,
			Mail:     email,
			Hash:     hashkey,
			Param:    "forgot",
			Created:  int(time.Now().Unix()),
			Livetime: int(time.Now().Unix()) + 172800,
		}

		db.Conn.Create(&timehash)

		//send email

		localizer := i18n.NewLocalizer(bundle, lg)

		subj := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "Subject",
			},
		})

		title := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "PassEmailTitle",
			},
			TemplateData: map[string]string{
				"Name": userExist.Name,
			},
		})

		keystring := localizer.MustLocalize(&i18n.LocalizeConfig{

			DefaultMessage: &i18n.Message{
				ID: "Key",
			},
			TemplateData: map[string]string{
				"Password": hashkey,
			},
		})

		linkarray := []string{keystring}
		ms := &sf.MailSettings{}
		ms.Custom = false

		go sf.SendHTMLEmail(userExist.Mail, title, linkarray, subj, "email.html", nil, ms)

		ans["response"] = "100201" // sent email with confirmation code
		ans["email"] = userExist.Mail

		return ans, nil, t

	}

	//check for key

	var userHash sf.TimeHash

	rows = db.Conn.Where(`hash = ? and mail = ?`, hashedkey, email).First(&userHash)
	if rows.RowsAffected == 0 {
		// return error. Hash is not exist in db

		ans["httpcode"] = 400
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000008",
			Message: "Hash is not exist in db",
		}
		errormsg = append(errormsg, errorans)

		return ans, errormsg, t
	}

	// delete hash
	db.Conn.Delete(sf.TimeHash{}, "hash = ? and mail = ?", hashedkey, email)
	// Create a new password
	userpass := sf.RandomString(12)

	//2.1 generete pass passhash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpass), 8)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("forgot.go: " + err.Error())
		}
	}

	//6. Write data to users table

	db.Conn.Table("users").Where("mail = ?", email).Updates(map[string]interface{}{"pass": hashedPassword})

	localizer := i18n.NewLocalizer(bundle, lg)

	subj := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Subject2",
		},
	})

	title := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PassEmailTitle",
		},
		TemplateData: map[string]string{
			"Name": userExist.Name,
		},
	})

	password := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "NewPass",
		},
		TemplateData: map[string]string{
			"Password": userpass,
		},
	})

	linkarray := []string{password}

	ms := &sf.MailSettings{}
	ms.Custom = false
	go sf.SendHTMLEmail(userExist.Mail, title, linkarray, subj, "email.html", nil, ms)
	//return data

	ans["response"] = "100202" // Password changed and sent email

	return ans, nil, t
}
