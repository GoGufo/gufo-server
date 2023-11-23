package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-server/gufodao"
	pb "github.com/gogufo/gufo-server/proto"
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
		Dbversion:  &t.Dbversion,
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
	}

	ans := sf.ToMapStringInterface(response.Data)

	moduleAnswerv3(w, r, ans, t)

}
