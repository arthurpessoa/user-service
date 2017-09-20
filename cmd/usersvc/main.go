package main

import (
	"flag"
	"os"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/arthurpessoa/user-service/internal/usersvc/validation"
	"github.com/arthurpessoa/user-service/internal/usersvc/user"
	"net/http"
	"os/signal"
	"syscall"
	"fmt"
)

const (
	defaultPort              = "8080"
)


func main() {
	var (
		addr  = envString("PORT", defaultPort)
		httpAddr          = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)


	//TODO: Dependency injection?
	users := user.NewUserRepository()

	var vs validation.Service
	vs = validation.NewService(users)


	httpLogger := log.With(logger, "component", "http")
	mux := http.NewServeMux()

	mux.Handle("/api/users/validate", validation.MakeHandler(vs, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())


	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}


func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}


func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
