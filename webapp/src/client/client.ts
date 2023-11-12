// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import * as crypto from 'crypto';


import {PluginSettings} from '../actions';
import {doPost} from '.'



export default class Client {
    
    public async startMeeting(settings:PluginSettings) :Promise<any>{
    
        const nonce = Math.round(Date.now() / 1000);
        const path = '/api/token';
        const myRoomId = settings.prefix + ':TheRoom';
        const body =
        {
            roomId: myRoomId,
            roomName: 'The room nname',
            userId: settings.user_id,
            role: 'guest',
            useLobby: false,
            estimatedNumberOfParticipants: 0,
            invalidateOnProperLeave: false,
            userData: {
                name: 'string',
            },
            roomLock: {
                enabled: false,
                isActive: false,
                mayBypass: false,
            },
            clientOptions: {
                disableCamera: false,
                disableChat: false,
                disableLayout: false,
                disableMembersList: false,
                disableMicrophone: false,
                disableMirroredCamera: false,
                disableRaiseHand: false,
                disableRecording: false,
                disableShareScreen: false,
                disableVisualEffects: false,
                forceVirtualBackgroundUrl: '',
                startAsMuted: false,
                startVideoFit: 'contain',
            },
        };
        const data = JSON.stringify(body);
        const hash:any = crypto.
            createHash('sha256').
            update(`${path}:${data}:${nonce}`).
            digest('hex');

        const buf = Buffer.from(settings.api_secret, 'base64');

        const key = crypto.createPrivateKey({
            key: buf,
        });

        const token = crypto.sign(null, hash, key).toString('hex');

        const headers = {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
            'X-Vidicue-Nonce': nonce,
            'X-Vidicue-Key': settings.api_key,
        };

        const res = await doPost(`${settings.compodium_api_url}${path}`, data, headers);
        return res;
    };
}

