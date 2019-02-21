package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	infoLog, errLog := getLogger()

	storage := getStorage(errLog)
	defer storage.Close()

	app := newApp(infoLog, errLog, storage)

	server := &http.Server{
		ErrorLog:     errLog,
		Addr:         getConf("LISTEN_ADDR"),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  20 * time.Second,
		Handler:      app.compileHandler(),
	}

	if err := server.ListenAndServe(); err != nil {
		errLog.Panicln(err)
	}
}

func getLogger() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INF: ", log.LUTC|log.LstdFlags|log.Lmicroseconds)
	errLog := log.New(os.Stderr, "ERR: ", log.LUTC|log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	return infoLog, errLog
}
