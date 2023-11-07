// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"encoding/json"
	"net/http"

	"runtime/debug"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

const (
	routeAPI             = "/api/v2"
	routeInstancePath    = "/instance/{id}"
	routeAPISettingsInfo = "/settingsinfo"
)

const routePrefixInstance = "instance"

func makeAPIRoute(path string) string {
	return routeAPI + path
}

func (p *Plugin) initializeRouter() {
	p.router = mux.NewRouter()
	p.router.Use(p.withRecovery)

	//instanceRouter := p.router.PathPrefix(routeInstancePath).Subrouter()

	apiRouter := p.router.PathPrefix(routeAPI).Subrouter()
	// User APIs

	apiRouter.HandleFunc(routeAPISettingsInfo, p.checkAuth(p.handleResponse(p.httpGetSettingsInfo))).Methods(http.MethodGet)

}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func respondErr(w http.ResponseWriter, code int, err error) (int, error) {
	http.Error(w, err.Error(), code)
	return code, err
}

func respondJSON(w http.ResponseWriter, obj interface{}) (int, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to marshal response"))
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return http.StatusInternalServerError, errors.WithMessage(err, "failed to write response")
	}
	return http.StatusOK, nil
}

func (p *Plugin) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.client.Log.Warn("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) checkAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func (p *Plugin) handleResponse(fn func(w http.ResponseWriter, r *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := fn(w, r)

		p.logResponse(r, status, err)
	}
}

func (p *Plugin) logResponse(r *http.Request, status int, err error) {
	if status == 0 || status == http.StatusOK {
		return
	}
	if err != nil {
		p.client.Log.Warn("ERROR: ", "Status", strconv.Itoa(status), "Error", err.Error(), "Path", r.URL.Path, "Method", r.Method, "query", r.URL.Query().Encode())
	}

	if status != http.StatusOK {
		p.client.Log.Debug("unexpected plugin response", "Status", strconv.Itoa(status), "Path", r.URL.Path, "Method", r.Method, "query", r.URL.Query().Encode())
	}
}
