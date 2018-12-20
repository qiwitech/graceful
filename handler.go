package graceful

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/trace"

	"github.com/golang/glog"
)

type Handlerer interface {
	http.Handler
	Handle(string, http.Handler)
}

type lazyDumper []byte

func (d lazyDumper) String() string {
	return hex.Dump(d)
}

func NewHandler(
	codec Codec,
	makebuf func() interface{},
	process func(context.Context, interface{}) (interface{}, error)) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			ts := time.Now()
			tr := trace.New(req.URL.Path, req.URL.Path)

			var err error
			defer func() {
				if err != nil {
					tr.LazyPrintf("err: %v", err)
					tr.SetError()
				}
				tr.Finish()
			}()

			ctx := req.Context()
			select {
			case <-ctx.Done():
				Error(w, ctx.Err(), http.StatusInternalServerError)
				return
			default:
			}

			data, err := ioutil.ReadAll(req.Body)
			if err != nil {
				Error(w, err, http.StatusInternalServerError)
				return
			}
			tr.LazyPrintf("raw: %v", lazyDumper(data))

			if err = req.Body.Close(); err != nil {
				Error(w, err, http.StatusInternalServerError)
				return
			}

			args := makebuf()
			if err = codec.Unmarshal(data, args); err != nil {
				Error(w, err, http.StatusInternalServerError)
				return
			}

			tr.LazyPrintf("req: %v", args)

			resp, err := process(ctx, args)
			if err != nil {
				Error(w, err, http.StatusInternalServerError)
				return
			}

			data, err = codec.Marshal(resp)
			if err != nil {
				Error(w, err, http.StatusInternalServerError)
				return
			}

			tr.LazyPrintf("resp: %v", resp)

			w.Header().Set("Content-Type", codec.MIME())
			w.Header().Set("X-Timing", time.Since(ts).String())
			w.WriteHeader(http.StatusOK)

			if _, err = w.Write(data); err != nil {
				glog.Error(err)
				return
			}
		})

}
