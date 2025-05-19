export type QoSLevel = 0 | 1 | 2

export interface Topic {
    id: string
    name: string
    subscribed: boolean
}

export interface Message {
    id: string
    topic: string
    payload: string
    qos: QoSLevel
    timestamp: string
}

export interface Credentials {
    ip: string
    port: string
    clientId: string
}
