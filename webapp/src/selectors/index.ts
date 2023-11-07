// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createSelector} from 'reselect';

import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';

import {GlobalState} from 'mattermost-redux/types/store';

import {id as PluginId} from '../manifest';
import {Instance} from '../types/model';

const getPluginState = (state: { [x: string]: any; }) => state['plugins-' + PluginId] || {};

export const getPluginServerRoute = (state: GlobalState) => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substring(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + PluginId;
};

export const getCurrentUserLocale = createSelector(
    getCurrentUser,
    (user) => {
        let locale = 'en';
        if (user && user.locale) {
            locale = user.locale;
        }

        return locale;
    },
);

export const isConnectModalVisible = (state: any) => getPluginState(state).connectModalVisible;
export const isDisconnectModalVisible = (state: any) => getPluginState(state).disconnectModalVisible;

export const isCreateModalVisible = (state: any) => getPluginState(state).createModalVisible;

export const getCreateModal = (state: any) => getPluginState(state).createModal;

export const isAttachCommentToIssueModalVisible = (state: any) => getPluginState(state).attachCommentToIssueModalVisible;

export const getAttachCommentToIssueModalForPostId = (state: any) => getPluginState(state).attachCommentToIssueModalForPostId;

export const getChannelIdWithSettingsOpen = (state: any) => getPluginState(state).channelIdWithSettingsOpen;

export const getChannelSubscriptions = (state: any) => getPluginState(state).channelSubscriptions;

export const isUserConnected = (state: any) => getUserConnectedInstances(state).length > 0;

export const canUserConnect = (state: any) => getPluginState(state).userCanConnect;

export const getUserConnectedInstances = (state: any): Instance[] => {
    const installed = getPluginState(state).installedInstances as Instance[];
    const connected = getPluginState(state).userConnectedInstances as Instance[];
    if (!installed || !connected) {
        return [];
    }

    return connected.filter((instance1) => installed.find((instance2) => instance1.instance_id === instance2.instance_id));
};

export const getInstalledInstances = (state: any): Instance[] => getPluginState(state).installedInstances;
export const instanceIsInstalled = (state: any): boolean => getInstalledInstances(state).length > 0;

export const getDefaultUserInstanceID = (state: any) => getPluginState(state).defaultUserInstanceID;

export const getPluginSettings = (state: any) => getPluginState(state).pluginSettings;
