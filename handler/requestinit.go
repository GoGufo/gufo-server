package handler

import (
	"net/http"
	"strings"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func RequestInit(r *http.Request) *pb.Request {
	t := &pb.Request{}
	p := bluemonday.UGCPolicy()

	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	module := p.Sanitize(patharray[3])
	t.Module = &module
	t.Path = &path
	t.Method = &r.Method

	sgn := viper.GetString("server.sign")
	curip := sf.ReadUserIP(r)
	usagent := r.UserAgent()

	t.Sign = &sgn

	t.IP = &curip

	t.UserAgent = &usagent

	//Function in Plugin
	if pathlenth >= 5 {
		*t.Param = p.Sanitize(patharray[4])
	}

	//ID for function in plugin
	if pathlenth >= 6 {
		*t.ParamID = p.Sanitize(patharray[5])

	}

	return t
}
