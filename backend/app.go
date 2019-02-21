package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/payfazz/go-router/segment"

	middleware "github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
)

type app struct {
	infoLog, errLog      *log.Logger
	storage              *storage
	adminName, adminPass string
}

func newApp(infoLog, errLog *log.Logger, storage *storage) *app {
	return &app{
		infoLog:   infoLog,
		errLog:    errLog,
		storage:   storage,
		adminName: getConf("ADMIN_NAME"),
		adminPass: getConf("ADMIN_PASS"),
	}
}

func (app *app) compileHandler() http.HandlerFunc {
	return middleware.Compile(
		logger.New(
			logger.DefaultLogger(app.infoLog),
		),

		paniclogger.New(
			15,
			paniclogger.DefaultLogger(app.errLog),
		),

		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				next(w, r)
			}
		},

		path.H{
			"/":        app.helloHandler(),
			"/counter": app.counterHandler(),
			"/reset":   app.resetHandler(),
		}.C(),
	)
}

func (app *app) helloHandler() http.HandlerFunc {
	return middleware.Compile(
		segment.E,
		defhandler.ResponseCodeWithMessage(200, "Hello World"),
	)
}

func (app *app) counterHandler() http.HandlerFunc {
	return middleware.Compile(
		segment.E,
		func(w http.ResponseWriter, r *http.Request) {
			if err := app.storage.incCounter(); err != nil {
				app.errLog.Println(err)
				defhandler.StatusInternalServerError(w, r)
				return
			}

			counter, err := app.storage.getCounter()
			if err != nil {
				app.errLog.Println(err)
				defhandler.StatusInternalServerError(w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Counter int `json:"counter"`
			}{
				Counter: counter,
			})
		},
	)
}

func (app *app) resetHandler() http.HandlerFunc {
	oke := defhandler.ResponseCode(200)
	return middleware.Compile(
		segment.E,
		method.H{
			"OPTIONS": func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization")
			},
			"POST": middleware.Compile(
				basicAuth(app.adminName, app.adminPass),
				func(w http.ResponseWriter, r *http.Request) {
					if err := app.storage.resetCounter(); err != nil {
						app.errLog.Println(err)
						defhandler.StatusInternalServerError(w, r)
						return
					}
					oke(w, r)
				},
			),
		}.C(),
	)
}
