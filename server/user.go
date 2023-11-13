// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

//nolint:revive,unused,stylecheck
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

func (p *Plugin) httpGetSettingsInfo(w http.ResponseWriter, r *http.Request) (int, error) {
	config := p.getConfiguration()
	return respondJSON(w, struct {
		CompodiumAPIURL string `json:"compodium_api_url"`
		APIKey          string `json:"api_key"`
		APISecret       string `json:"api_secret"`
		Prefix          string `json:"prefix"`
		UserId          string `json:"user_id"`
	}{
		CompodiumAPIURL: config.CompodiumAPIURL,
		APIKey:          config.APIKey,
		APISecret:       config.APISecret,
		Prefix:          config.Prefix,
		UserId:          config.UserId,
	})
}

type startMeetingRequest struct {
	ChannelID  string `json:"channel_id"`
	MeetingURL string `json:"meeeting_url"`
}

func (p *Plugin) postMeeting(creatorUsername string, meetingURL string, channelID string) (*model.Post, *model.AppError) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf("Meeting started at %s.", meetingURL),
		Type:      "custom_compodium",
		Props: map[string]interface{}{
			"meeting_link":             meetingURL,
			"meeting_creator_username": creatorUsername,
		},
	}
	return p.API.CreatePost(post)
}

func (p *Plugin) handleStartMeeting(w http.ResponseWriter, r *http.Request) (int, error) {
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return http.StatusUnauthorized, nil
	}

	var req startMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusBadRequest, nil
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return appErr.StatusCode, nil
	}

	if _, appErr = p.API.GetChannelMember(req.ChannelID, userID); appErr != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return http.StatusForbidden, nil
	}

	createdPost, appErr := p.postMeeting(user.Username, req.MeetingURL, req.ChannelID)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return appErr.StatusCode, nil
	}

	if appErr = p.API.KVSet(fmt.Sprintf("%v%v", "", ""), []byte(createdPost.Id)); appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return appErr.StatusCode, nil
	}

	if _, err := w.Write([]byte(fmt.Sprintf(`{"meeting_url": "%s"}`, req.MeetingURL))); err != nil {
		p.API.LogWarn("failed to write response", "error", err.Error())
	}
	return http.StatusCreated, nil
}
