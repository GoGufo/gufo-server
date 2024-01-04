package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func connectgrpc(w http.ResponseWriter, r *http.Request, t *sf.Request) {

	pluginname := fmt.Sprintf("microservices.%s", t.Module)
	hostpath := fmt.Sprintf("%s.host", pluginname)
	portpath := fmt.Sprintf("%s.port", pluginname)
	host := viper.GetString(hostpath)
	port := viper.GetString(portpath)

	connection := fmt.Sprintf("%s:%s", host, port)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	//	args := os.Args
	conn, err := grpc.Dial(connection, opts...)

	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}

		nomoduleAnswerv3(w, r)
		return
	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)
	request := &pb.Request{
		Module:     &t.Module,
		Param:      &t.Param,
		ParamID:    &t.ParamID,
		Action:     &t.Action,
		Args:       sf.ToMapStringAny(t.Args),
		Token:      &t.Token,
		TokenType:  &t.TokenType,
		TimeStamp:  sf.Int32(t.TimeStamp),
		Language:   &t.Language,
		APIVersion: &t.APIVersion,
		Method:     &r.Method,
		UID:        &t.UID,
		IsAdmin:    sf.Int32(t.IsAdmin),
		SessionEnd: sf.Int32(t.SessionEnd),
		Completed:  sf.Int32(t.Completed),
		Readonly:   sf.Int32(t.Readonly),
	}

	response, err := client.Do(context.Background(), request)

	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}
		nomoduleAnswerv3(w, r)
		return
	}

	ans := sf.ToMapStringInterface(response.Data)

	//update *sf.Request
	if *response.RequestBack.Token != t.Token {
		t.Token = *response.RequestBack.Token
	}
	if *response.RequestBack.TokenType != t.TokenType {
		t.TokenType = *response.RequestBack.TokenType
	}
	if *response.RequestBack.Language != t.Language {
		t.Language = *response.RequestBack.Language
	}
	if *response.RequestBack.UID != t.UID {
		t.UID = *response.RequestBack.UID
	}
	if int(*response.RequestBack.IsAdmin) != t.IsAdmin {
		t.IsAdmin = int(*response.RequestBack.IsAdmin)
	}
	if int(*response.RequestBack.SessionEnd) != t.SessionEnd {
		t.SessionEnd = int(*response.RequestBack.SessionEnd)
	}
	if int(*response.RequestBack.Completed) != t.Completed {
		t.Completed = int(*response.RequestBack.Completed)
	}
	if int(*response.RequestBack.Readonly) != t.Readonly {
		t.Readonly = int(*response.RequestBack.Readonly)
	}

	moduleAnswerv3(w, r, ans, t)

}
