export const endpoints = {
    // connection
    credentials: "/credentials",
    disconnect: "/disconnect",
    ping: "/ping",

    // topics
    subscribeToTopic: "/topic/subscribe",
    unsubscribeToTopic: "/topic/unsubscribe",
    subscribedTopics: "/topic/subscribed",
    getAllTopics: "/topic/all-known-subscribed",

    // messages
    sendMessageToTopic: "/topic/send-message",
    getMessageFromTopic: "/topic/messages",
}