// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	/*
		"github.com/mattermost/mattermost/server/public/model"
		"github.com/mattermost/mattermost/server/public/plugin"
		"github.com/mattermost/mattermost/server/public/pluginapi"
	*/)

const (
	botUserName    = "compodium"
	botDisplayName = "Compodium"
	botDescription = "Created by the Compodium plugin."

	trueString  = "true"
	falseString = "false"

	//zoomProviderName = "Zoom"
)

type Plugin struct {
	plugin.MattermostPlugin
	client *pluginapi.Client

	router *mux.Router

	conf configuration

	// botUserID of the created bot account.
	botUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	siteURL string

	//telemetryClient telemetry.Client
	//tracker         telemetry.Tracker
}

// OnActivate checks if the configurations is valid and ensures the bot account exists
func (p *Plugin) OnActivate() error {

	p.client = pluginapi.NewClient(p.API, p.Driver)

	p.client.System.GetBundlePath()

	if err := p.registerSiteURL(); err != nil {
		return errors.Wrap(err, "could not register site URL")
	}
	/*
		command, err := p.getCommand()
		if err != nil {
			return errors.Wrap(err, "failed to get command")
		}


		err = p.API.RegisterCommand(command)
		if err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	*/

	botUserID, err := p.client.Bot.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: botDescription,
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}
	p.botUserID = botUserID

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "couldn't get bundle path")
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "profile.png"))
	if err != nil {
		return errors.Wrap(err, "couldn't read profile image")
	}

	if appErr := p.API.SetProfileImage(botUserID, profileImage); appErr != nil {
		return errors.Wrap(appErr, "couldn't set profile image")
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {
	return nil
}

// registerSiteURL fetches the site URL and sets it in the plugin object.
func (p *Plugin) registerSiteURL() error {
	siteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if siteURL == nil || *siteURL == "" {
		return errors.New("could not fetch siteURL")
	}

	p.siteURL = *siteURL
	return nil
}
