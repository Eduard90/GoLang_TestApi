package middlewares

import (
	"net/http"
	"fmt"
)

func TestMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Before")
	next(rw, r)
	fmt.Println("After")
}