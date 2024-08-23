package handlers

import "net/http"

func GetRootHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "static/index.html")
}
