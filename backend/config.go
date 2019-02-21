package main

import (
	"os"
)

var defConfig = map[string]string{
	"DATABASE_URI": "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	"LISTEN_ADDR":  ":8080",
	"ADMIN_NAME":   "admin",
	"ADMIN_PASS":   "admin",
}

func getConf(key string) string {
	res := os.Getenv(key)
	if res == "" {
		res = defConfig[key]
	}
	return res
}
