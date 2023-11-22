package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-server/gufodao"
	pb "github.com/gogufo/gufo-server/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func connectgrpc(w http.ResponseWriter, r *http.Request, mod string, t *sf.Request) {

	pluginname := fmt.Sprintf("microservices.%s", t.Module)
	host := fmt.Sprintf("%s.host", pluginname)
	port := fmt.Sprintf("%s.port", pluginname)
	connection := fmt.Sprintf("%s:%s", host, port)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	args := os.Args
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
		Message: args[1],
	}
	response, err := client.Do(context.Background(), request)

	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}
	}

	fmt.Println(response.Message)

}
