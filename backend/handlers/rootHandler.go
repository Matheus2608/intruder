package handlers

import "net/http"

func GetRootHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "frontend/templates/index.html")
}
