export type QoSLevel = 0 | 1 | 2

export interface Topic {
    BrokerId: string;
    UserId: string;
    Id: string;
    Topic: string;
    CreationDate: string;
    Subscribed: boolean
}

export interface Message {
    id: string
    topic: string
    payload: string
    qos: QoSLevel
    timestamp: string
    ClientId: string
}

export interface Credentials {
    ip: string
    port: string
    clientId: string
}