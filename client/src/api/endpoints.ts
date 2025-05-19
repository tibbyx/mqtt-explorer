/*
* Defines the API endpoints as constants
* It`s easier to update URLs if we have them in one place
*/

export const endpoints = {
    // connection
    credentials: "/credentials",
    disconnect: "/disconnect",
    ping: "/ping",

    // topics
    subscribeToTopic: "/topic/subscribe",
    unsubscribeToTopic: "/topic/unsubscribe",
    subscribedTopics: "/topic/subscribed",

    // messages
    sendMessageToTopic: "/topic/send-message",
    getMessageFromTopic: "/topic/messages",
}