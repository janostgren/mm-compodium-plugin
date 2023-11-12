import { GlobalState } from 'mattermost-redux/types/store';

import { getPluginServerRoute } from '../selectors';
import ActionTypes from '../action_types';
import {PostTypes} from 'mattermost-redux/action_types';
import { doFetch,doPost } from '../client';
import Client from '../client/client'

export interface PluginSettings {
    compodium_api_url: string,
    api_key: string,
    api_secret: string
    prefix: string
    user_id: string
}

const client:Client =new Client()

export function getSettings():PluginSettings {
    const settings:any=async (dispatch: (arg0: { type: any; data: any; }) => void, getState: () => GlobalState) => {
        let data;
        const baseUrl = getPluginServerRoute(getState());
        try {
            data = await doFetch(`${baseUrl}/api/v2/settingsinfo`, {
                method: 'get',
            });

            dispatch({
                type: ActionTypes.RECEIVED_PLUGIN_SETTINGS,
                data,
            });
        } catch (error) {
            return {error};
        }

        return data;
    };
    return <PluginSettings> settings
    
}

export function startMeeting() {
    return async (dispatch: (arg0: { type: any; data: any;channelId:string }) => void, getState: () => GlobalState) => {
        let data;
        const baseUrl = getPluginServerRoute(getState());
        const channelId= getState().entities.channels.currentChannelId
        try {

            const meetingInfo:any = await client.startMeeting(getSettings())
            if( meetingInfo.url) {
                window.open(meetingInfo.url);
            }

            data = await doPost(`${baseUrl}/api/v2/startMeeting`, {MeetingUrl: meetingInfo.url,ChannelID: channelId});
           

        } catch (error:any) {
            let m = error.message;
            if (error.message && error.message[0] === '{') {
                const e = JSON.parse(error.message);

                // Error is from Compodium API
                if (e && e.message) {
                    m = '\nCompodium error: ' + e.message;
                }
            }
            const post = {
                id: 'compodiumPlugin' + Date.now(),
                create_at: Date.now(),
                update_at: 0,
                edit_at: 0,
                delete_at: 0,
                is_pinned: false,
                user_id: getState().entities.users.currentUserId,
                channel_id: channelId,
                root_id: '',
                parent_id: '',
                original_id: '',
                message: m,
                type: 'system_ephemeral',
                props: {},
                hashtags: '',
                pending_post_id: '',
            };

            dispatch({
                type: PostTypes.RECEIVED_NEW_POST,
                data: post,
                channelId: channelId
        
            });
           
        
            return {error};
        }

        return data;
    };
}