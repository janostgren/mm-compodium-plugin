import React from 'react';

import {id as pluginId} from './manifest';

import ChannelHeaderIcon from './components/channel-header-icon/channel-header-icon';
import {startMeeting} from './actions';

class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry:any, store:any) {
        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderIcon/>,
            (channel:any) => {
                startMeeting();
            },
            'Start Compodium Meeting',
        );

        //registry.registerPostTypeComponent('custom_compodium', PostTypeBluejeans);
    }
}
declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(pluginId, new Plugin());