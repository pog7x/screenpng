package server

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"net/http"
	"os"

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
	s.registerHandlers()
	s.logger.Sugar().Infof("Server is running on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s HTTPServer) registerHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	mux.Handle("/screenshot", http.HandlerFunc(s.screenshot()))
	s.server.Handler = mux
}

func (s HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s HTTPServer) screenshot() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		name := r.URL.Query().Get("name")

		handler := func(screenshotBytes []byte) error {
			screenshot, _, err := image.Decode(bytes.NewReader(screenshotBytes))
			if err != nil {
				s.logger.Error("Decoding screenshot bytes error", zap.Error(err))
				return err
			}

			out, err := os.Create(name)
			if err != nil {
				s.logger.Sugar().Errorf("Creating screenshot file %s error %v", "sd", err)
				return err
			}

			err = png.Encode(out, screenshot)
			if err != nil {
				s.logger.Error("Encoding screenshot bytes error", zap.Error(err))
				return err
			}
			return nil
		}

		var maximize string
		go s.f.MakeScreenshot(
			ssfactory.MakeScreenshotPayload{
				URL:            url,
				DOMElementBy:   ssfactory.ByTagName,
				DOMElementName: "body",
				Scroll:         true,
				BytesHandler:   handler,
				MaximizeWindow: &maximize,
			},
		)
		w.Write([]byte{})
	}
}
