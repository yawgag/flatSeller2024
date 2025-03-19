package mainPage

import "net/http"

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../pkg/template/mainPage.html")
}
