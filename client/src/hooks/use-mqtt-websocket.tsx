import {useCallback, useEffect, useState} from "react"
import {v4 as uuidv4} from "uuid"
import type {Message, QoSLevel, Topic} from "../lib/types"

interface MqttWebSocketOptions {
    onConnect?: () => void
    onDisconnect?: () => void
    searchQuery?: string
    shouldConnect?: boolean
}

export function useMqttWebSocket({onConnect,onDisconnect, searchQuery = "",shouldConnect = false}: MqttWebSocketOptions = {}) {
    const [connected, setConnected] = useState(false)
    const [topics, setTopics] = useState<Topic[]>([])
    const [messages, setMessages] = useState<Message[]>([])
    //const reconnectTimeoutRef = useRef<number | null>(null)

    // Filter topics
    const filteredTopics = topics.filter((topic) => topic.name.toLowerCase().includes(searchQuery.toLowerCase()))

    // Mock connection to simulate connection to Go-Fiber
    const connectWebSocket = useCallback(() => {

        // Simulate successful connection after a delay
        setTimeout(() => {
            setConnected(true)
            onConnect?.()

            // Load some sample topics
            if (topics.length === 0) {
                setTopics([
                    {id: uuidv4(), name: "Brot", subscribed: true},
                    {id: uuidv4(), name: "Golf", subscribed: true},
                    {id: uuidv4(), name: "Espresso", subscribed: true},
                    {id: uuidv4(), name: "Schule", subscribed: false},
                    {id: uuidv4(), name: "Arbeit", subscribed: false},
                ])
            }
        }, 1000)
    }, [onConnect, topics.length])

    const disconnectWebSocket = useCallback(() => {
        setConnected(false)
        onDisconnect?.()
    }, [onDisconnect])

    // Connect on mount
    useEffect(() => {
        if (shouldConnect && !connected) {
            connectWebSocket()
        } else if (!shouldConnect && connected) {
            disconnectWebSocket()
        }
    }, [shouldConnect, connected, connectWebSocket, disconnectWebSocket])

    // Create a new topic
    const createTopic = useCallback((name: string) => {
        const newTopic: Topic = {
            id: uuidv4(),
            name,
            subscribed: false,
        }

        setTopics((prev) => [...prev, newTopic])
        return newTopic
    }, [])

    // Delete a topic
    const deleteTopic = useCallback((id: string) => {
        setTopics((prev) => prev.filter((topic) => topic.id !== id))
    }, [])

    // Rename a topic
    const renameTopic = useCallback((id: string, newName: string) => {
        setTopics((prev) => prev.map((topic) => (topic.id === id ? {...topic, name: newName} : topic)))
    }, [])

    // Subscribe to a topic
    const subscribeTopic = useCallback((id: string) => {
        setTopics((prev) => {
            return prev.map((topic) => (topic.id === id ? {...topic, subscribed: true} : topic))
        })
    }, [])

    // Unsubscribe from a topic
    const unsubscribeTopic = useCallback((id: string) => {
        setTopics((prev) => {
            return prev.map((topic) => (topic.id === id ? {...topic, subscribed: false} : topic))
        })
    }, [])

    // Publish a message
    const publishMessage = useCallback((topic: string, payload: string, qos: QoSLevel = 0) => {
        const newMessage: Message = {
            id: uuidv4(),
            topic,
            payload,
            qos,
            timestamp: new Date().toISOString(),
        }

        // Add the message to our state
        setMessages((prev) => [...prev, newMessage])

        return newMessage
    }, [])


    return {
        connect: connectWebSocket,
        disconnect: disconnectWebSocket,
        connected,
        topics,
        messages,
        filteredTopics,
        createTopic,
        deleteTopic,
        renameTopic,
        subscribeTopic,
        unsubscribeTopic,
        publishMessage,
    }
}
