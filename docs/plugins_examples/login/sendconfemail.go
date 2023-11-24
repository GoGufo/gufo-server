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

	"github.com/BurntSushi/toml"
	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func SendConfEmail(email string, uid string, uname string, lang string) interface{} {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	sysdir := viper.GetString("server.sysdir")
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/en.toml", sysdir))
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/lang/auth/ru.toml", sysdir))

	//Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog(err.Error())
		}

		return ""
	}

	hash := sf.Stringen(64)

	timehash := sf.TimeHash{
		UID:      uid,
		Mail:     email,
		Hash:     hash,
		Param:    "signup",
		Created:  int(time.Now().Unix()),
		Livetime: int(time.Now().Unix()) + 172800,
	}

	db.Conn.Create(&timehash)

	lg := ""
	switch lang {
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

	//4. Send Confirmation email
	domain := viper.GetString("server.domain")
	link := fmt.Sprintf("%s/confirmemail?email=%s&token=%s&lang=%s", domain, email, hash, lg)
	subj := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "ConfSubj",
		},
	})

	title := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "ConfirmEmailTitle",
		},
		TemplateData: map[string]string{
			"Name": uname,
		},
	})

	conflinka := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "ConfirmEmailLinka",
		},
		TemplateData: map[string]string{
			"URL": link,
		},
	})

	conflinkb := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "ConfirmEmailLinkb",
		},
		TemplateData: map[string]string{
			"URL": link,
		},
	})

	linkarray := []string{conflinka, conflinkb}

	ms := &sf.MailSettings{}
	ms.Custom = false
	go sf.SendHTMLEmail(email, title, linkarray, subj, "email.html", nil, ms)

	return ""
}
