// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';


export const doFetch = async (url:string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

   
    if (response.ok) {
        return response.json();        
    }

    const data = await response.text();

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export const doPost = async (url:string , body:any, headers = {}) => {
    const options = {
        method: 'post',
        body,
        headers,
    };

    return await doFetch(url,options)

};

