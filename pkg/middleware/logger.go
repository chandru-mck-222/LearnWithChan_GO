package middleware

import "log"

func Condlog(format string, v ...any) {
	log.Printf(format, v...)
}
