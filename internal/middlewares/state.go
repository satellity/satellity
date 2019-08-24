package middlewares

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime"
	"satellity/internal/configs"
	"satellity/internal/session"
	"satellity/internal/views"
	"time"
)

// State output states of request, e.g.: r.Method, r.URL etc.
func State(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		session.Logger(r.Context()).Infof("INFO -- : Started %s '%s'", r.Method, r.URL)
		defer func() {
			session.Logger(r.Context()).Infof("INFO -- : Completed %s in %fms", r.Method, time.Now().Sub(start).Seconds())
		}()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
			return
		}
		if len(body) > 0 {
			session.Logger(r.Context()).Infof("INFO -- : Paremeters %s", string(body))
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		r = r.WithContext(session.WithRequestBody(r.Context(), string(body)))
		w.Header().Set("X-Build-Info", configs.BuildVersion+"-"+runtime.Version())
		handler.ServeHTTP(w, r)
	})
}
