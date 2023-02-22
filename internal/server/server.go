package server

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pog7x/screenpng/configs"
	"github.com/pog7x/ssfactory"

	"net/http/pprof"

	"go.uber.org/zap"
)

type HTTPServer struct {
	logger *zap.Logger
	server *http.Server
	f      ssfactory.Factory
}

func NewHTTPServer(logger *zap.Logger, cfg *configs.Config, f ssfactory.Factory) HTTPServer {
	return HTTPServer{
		logger: logger,
		f:      f,
		server: &http.Server{
			Addr:         cfg.ServerListenAddr,
			ReadTimeout:  cfg.ServerReadTimeout,
			WriteTimeout: cfg.ServerWriteTimeout,
		},
	}
}

func (s HTTPServer) Start() error {
	s.server.Handler = s.registerRouter()
	s.logger.Sugar().Infof("Server is running on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s HTTPServer) registerRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(s.loggingMiddleware)

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	r.HandleFunc("/screenshot", screenshot(s.logger, s.f)).Methods("POST")

	return r
}

func (s HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func screenshot(logger *zap.Logger, factory ssfactory.Factory) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type screenshotReq struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		var screenshotReqBody screenshotReq
		err = json.Unmarshal(body, &screenshotReqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		handler := func(screenshotBytes []byte) error {
			screenshot, _, err := image.Decode(bytes.NewReader(screenshotBytes))
			if err != nil {
				logger.Error("Decoding screenshot bytes error", zap.Error(err))
				return err
			}

			out, err := os.Create(screenshotReqBody.Name)
			if err != nil {
				logger.Sugar().Errorf("Creating screenshot file %s error %v", "sd", err)
				return err
			}

			err = png.Encode(out, screenshot)
			if err != nil {
				logger.Error("Encoding screenshot bytes error", zap.Error(err))
				return err
			}
			return nil
		}

		var maximize string
		go factory.MakeScreenshot(
			ssfactory.MakeScreenshotPayload{
				URL:            screenshotReqBody.URL,
				DOMElementBy:   ssfactory.ByTagName,
				DOMElementName: "body",
				Scroll:         true,
				BytesHandler:   handler,
				MaximizeWindow: &maximize,
			},
		)
		w.WriteHeader(http.StatusOK)
	}
}

func (s HTTPServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("New request", zap.String("request_uri", r.RequestURI), zap.String("method", r.Method))
		next.ServeHTTP(w, r)
	})
}
