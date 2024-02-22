package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

		errorAnswer(w, r, t, 400, "0000234", err.Error())

	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)
	request := &pb.Request{
		Module:     &t.Module,
		Param:      &t.Param,
		ParamID:    &t.ParamID,
		Path:       &t.Path,
		Action:     &t.Action,
		Args:       sf.ToMapStringAny(t.Args),
		Token:      &t.Token,
		Sign:       &t.Sign,
		IP:         &t.IP,
		UserAgent:  &t.UserAgent,
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

	if r.Method == "PUT" {
		/*
			var (
				buf        []byte
				firstChunk bool
			)
		*/
		//PUT mean file upload, so we check for file data
		file, handler, err := r.FormFile("file")

		if err != nil || file == nil || handler.Filename == "" {
			errorAnswer(w, r, t, 400, "0000235", "Missing File")

		}

		defer file.Close()

		buft := bytes.NewBuffer(nil)
		if _, err := io.Copy(buft, file); err != nil {

			errorAnswer(w, r, t, 400, "0000235", err.Error())
		}

		request.Filename = &handler.Filename
		request.File = buft.Bytes()

		/*
			//start uploader
			buf = make([]byte, chunkSize)
			firstChunk = true
			for {
				n, errRead := file.Read(buf)
				if errRead != nil {
					if errRead == io.EOF {
						errRead = nil
						break
					}
					errorAnswer(w, r, t, 400, "0000235", "errored while copying from file to buf")
				}

				if firstChunk {
					request.Filename = &handler.Filename
					request.File = buf[:n]
					firstChunk = false
				} else {
					request.File = buf[:n]
				}
				if err != nil {
					errorAnswer(w, r, t, 400, "0000235", "failed to send chunk via stream file")
				}

			}
		*/
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

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}
