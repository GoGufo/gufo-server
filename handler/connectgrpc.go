package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	pbv "gopkg.in/cheggaaa/pb.v1"
)

const chunkSize = 64 * 1024

type uploader struct {
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	pool        *pbv.Pool
	DoneRequest chan string
	FailRequest chan string
}

func connectgrpc(w http.ResponseWriter, r *http.Request, t *pb.Request) {

	pluginname := fmt.Sprintf("microservices.%s", *t.Module)
	hostpath := fmt.Sprintf("%s.host", pluginname)
	portpath := fmt.Sprintf("%s.port", pluginname)
	host := viper.GetString(hostpath)
	port := viper.GetString(portpath)

	plygintype := fmt.Sprintf("%s.type", pluginname)

	if plygintype == "internal" {

		if len(r.Header["X-Sign"]) == 0 {
			errorAnswer(w, r, t, 401, "0000234", "You have no rights")
			return
		}
		//Check for X-Sign
		signheader := r.Header["X-Sign"][0]
		sgn := viper.GetString("server.sign")

		if sgn != signheader {
			errorAnswer(w, r, t, 401, "0000234", "You have no rights")
			return
		}

	}

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

		errorAnswer(w, r, t, 400, "0000234", err.Error())
		return

	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)

	response, err := client.Do(context.Background(), t)

	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("connectgrpc: " + err.Error())
		}
		errorAnswer(w, r, t, 500, "0000236", fmt.Sprintf("Module connection error: %s", err.Error()))
		return
	}

	ans := sf.ToMapStringInterface(response.Data)

	//update *sf.Request
	if response.RequestBack.Token != t.Token {
		t.Token = response.RequestBack.Token
	}
	if response.RequestBack.TokenType != t.TokenType {
		t.TokenType = response.RequestBack.TokenType
	}
	if response.RequestBack.Language != t.Language {
		t.Language = response.RequestBack.Language
	}
	if response.RequestBack.UID != t.UID {
		t.UID = response.RequestBack.UID
	}
	if response.RequestBack.IsAdmin != t.IsAdmin {
		t.IsAdmin = response.RequestBack.IsAdmin
	}
	if response.RequestBack.SessionEnd != t.SessionEnd {
		t.SessionEnd = response.RequestBack.SessionEnd
	}
	if response.RequestBack.Completed != t.Completed {
		t.Completed = response.RequestBack.Completed
	}
	if response.RequestBack.Readonly != t.Readonly {
		t.Readonly = response.RequestBack.Readonly
	}

	moduleAnswerv3(w, r, ans, t)

}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}
