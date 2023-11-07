
export enum InstanceType {
    CLOUD = 'cloud',
    SERVER = 'server',
}

export type Instance = {
    alias?: string;
    instance_id: string;
    type: InstanceType;
}

export type GetConnectedResponse = {
    data: {
        can_connect: boolean;
        instances: Instance[];
        is_connected: boolean;
        user: {
            connected_instances: Instance[];
            default_instance_id?: string;
        };
    };
    error?: Error;
};

export type APIResponse<T> = {
    error?: Error;
    data: T;
};

