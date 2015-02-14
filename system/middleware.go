package system

import (
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gorilla/context"
	"github.com/zenazn/goji/web"
)

func JSON(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctAll := r.Header.Get("Content-Type")
		ct := strings.Split(ctAll, ";")
		if ct[0] == "application/json" {
			var params Params
			c.Env["ParamsError"] = json.NewDecoder(r.Body).Decode(&params)
			c.Env["Params"] = params
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) Session(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := application.Store.Get(r, "session")
		c.Env["Session"] = session
		h.ServeHTTP(w, r)
		context.Clear(r)
	}
	return http.HandlerFunc(fn)
}
