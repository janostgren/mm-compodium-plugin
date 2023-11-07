// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"net/http"
)

func (p *Plugin) httpGetSettingsInfo(w http.ResponseWriter, r *http.Request) (int, error) {
	config := p.getConfiguration()
	return respondJSON(w, struct {
		CompodiumAPIURL string `json:"compodium_api_url"`
		APIKey          string `json:"api_key"`
		APISecret       string `json:"api_secret"`
		Prefix          string `json:"prefix"`
	}{
		CompodiumAPIURL: config.CompodiumAPIURL,
		APIKey:          config.APIKey,
		APISecret:       config.APISecret,
		Prefix:          config.Prefix,
	})
}
