# Gufo

## General Information
Gufo is a basement API framework server. With Gufo you can create any API server you want. Just need to write a plugin with your features and connect it to Gufo.

Gufo API v3 is RESTfull. It supports next methods:
- GET
- POST
- OPTIONS
- DELETE
- PUT

Currently Gufo supports Plugins only. But our target is create 100% microservices architecture.

## How Create Plugin

Gufo is 100% ready to go. No need to update any codes.
Gufo is working with Go Plugins library.
API url structure is {hostname}/api/{api_version}/{plugin_name}/{plugin_function}/{any_other_data_requested_by_your_plugin}

1. Create Function with name "Init" and function "main"
```
package main

import (
  sf "github.com/gogufo/gufodao" //This is important library that content all necessary functions for communications with Gufo Server
  )

const VERSIONPLUGIN = "1.0" //Check for Plugin version
const VERSIONDB = "1.0" //Check for DB version

func main() {
  //If you use Sentry server. Connection to Sentry should be in this function
}

func Init(t *sf.Request, r *http.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {

  ans := make(map[string]interface{})
  var errormsg []sf.ErrorMsg

//Some of your codes and connections to other features

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

//If your need access to authorised users only, please add next check. If your plugin only for unauthorised user, use "!=" instead of "="
if t.UID == "" {
  ans["httpcode"] = 401
  errormsg := []sf.ErrorMsg{}
  errorans := sf.ErrorMsg{
    Code:    "000011",
    Message: "You are not authorised",
  }
  errormsg = append(errormsg, errorans)
  return ans, errormsg, t
}

//This is switch for connect to necessary functions depends from URL
//API url structure is {hostname}/api/{api_version}/{plugin_name}/{t.Param}...

switch t.Param {
case "something":
  ans, errormsg, t = Something(t, r)
case "info":
  ans, errormsg, t = info(t)
default:
  ans["httpcode"] = 404
  errormsg := []sf.ErrorMsg{}
  errorans := sf.ErrorMsg{
    Code:    "000012",
    Message: "Missing argument",
  }
  errormsg = append(errormsg, errorans)
  return ans, errormsg, t

}


  return ans, errormsg, t
}

func info(t *sf.Request) (map[string]interface{}, []sf.ErrorMsg, *sf.Request) {
  //We suggest to add info function to provide information about plugin
	ans := make(map[string]interface{})
	ans["pluginname"] = "Info about plugin"
	ans["version"] = VERSIONPLUGIN
	ans["apiversion"] = VERSIONDB
	ans["description"] = "This plugin for something you want
	return ans, nil, t
}

```

2. For build plugin use

```
go build -buildmode=plugin -o /var/lib/{plugin_name}.so  plugins/{your_plugin_name}}/*.go

```

Go plugins are work very strictly. If some dependencies between Gufo Server and Plugin are different - it will return error. TO avoid this problem we suggest to build Gufo Server and Plugin together.

#Ho to update plugin
You may update plugins. Just need to replace old .so file with new one. After replacement need to restart Gufo Server
