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
	"net/http"
	"time"

	"github.com/certifi/gocertifi"
	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-server/gufodao"
	"github.com/spf13/viper"
)

const VERSIONPLUGIN = "1.0"
const VERSIONDB = "1.0"

func main() {

	if viper.GetBool("server.sentry") {

		sf.SetLog("Connect to Setry...")

		sentryClientOptions := sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			EnableTracing:    viper.GetBool("sentry.tracing"),
			Debug:            viper.GetBool("sentry.debug"),
			TracesSampleRate: viper.GetFloat64("sentry.trace"),
		}

		rootCAs, err := gocertifi.CACerts()
		if err != nil {
			sf.SetLog("Could not load CA Certificates for Sentry: " + err.Error())

		} else {
			sentryClientOptions.CaCerts = rootCAs
		}

		err = sentry.Init(sentryClientOptions)

		if err != nil {
			sf.SetLog("Error with sentry.Init: " + err.Error())
		}

		flushsec := viper.GetDuration("sentry.flush")

		defer sentry.Flush(flushsec * time.Second)

	}

}

func Init(t *sf.Request, r *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

	ans := make(map[string]interface{})
	var errormsg []sf.ErrorMsg

	// Check if Plugin DB Version is same as App DB Version
	if t.Dbversion != VERSIONDB {
		//DB version is missmuched
		ans["httpcode"] = 409
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000006",
			Message: "DB version is missmuched",
		}
		errormsg = append(errormsg, errorans)
		return ans, errormsg, t
	}

	if r.Method == "POST" {

		if t.Token != "" {
			ans["httpcode"] = 400
			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000010",
				Message: "You are already has session",
			}
			errormsg = append(errormsg, errorans)
			return ans, errormsg, t
		}

	}

	switch t.Param {
	case "signin":
		ans, errormsg, t = Signin(t, r)
	case "forgot":
		ans, errormsg, t = Forgot(t)
	case "confemail":
		ans, errormsg, t = confirmemail(t)
	case "otp":
		ans, errormsg, t = otp(t, r)
	default:

		ans["httpcode"] = 404
		errormsg := []sf.ErrorMsg{}
		errorans := sf.ErrorMsg{
			Code:    "000012",
			Message: "Missing argument",
		}
		errormsg = append(errormsg, errorans)

	}

	return ans, errormsg, t
}
