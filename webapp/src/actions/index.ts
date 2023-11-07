import {GlobalState} from 'mattermost-redux/types/store';

import {getPluginServerRoute} from '../selectors';
import ActionTypes from '../action_types';
import {doFetch} from '../client';

export function getSettings() {
    return async (dispatch: (arg0: { type: any; data: any; }) => void, getState: () => GlobalState) => {
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
}