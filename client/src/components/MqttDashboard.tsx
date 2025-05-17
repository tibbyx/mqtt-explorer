import {Header} from "@/components/Header.tsx";
import TopicPanel from "@/components/topic/TopicPanel.tsx";
import {MessagePanel} from "@/components/message/MessagePanel.tsx";
import {ConnectionPanel} from "@/components/connection/ConnectionPanel.tsx";
import type {Topic} from "@/lib/types.ts";
import {useMqttWebSocket} from "@/hooks/use-mqtt-websocket.tsx";
import {useToast} from "../hooks/use-toast"
import {useState} from "react";

function MqttDashboard() {
    const [selectedTopic, setSelectedTopic] = useState<Topic | null>(null)
    const [searchQuery, setSearchQuery] = useState("")
    const [isConnected, setIsConnected] = useState(false);
    const {toast} = useToast()
    const {
        messages,
        topics,
        createTopic,
        deleteTopic,
        renameTopic,
        subscribeTopic,
        unsubscribeTopic,
        filteredTopics,
        publishMessage
    } = useMqttWebSocket({
        onConnect: () => {
            toast({
                title: "Connected",
                description: "Successfully connected to MQTT broker",
                className: "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-900/50",
            })
        },
        onDisconnect: () => {
            toast({
                title: "Disconnected",
                description: "Connection to MQTT broker lost",
                variant: "destructive",
            })
        },
        searchQuery,
    })

    // Handle topic subscription
    const handleSubscribe = (topicId: string) => {
        subscribeTopic(topicId)
        // Update the selected topic if it's the one being subscribed to
        if (selectedTopic && selectedTopic.id === topicId) {
            const updatedTopic = topics.find((t: Topic) => t.id === topicId)
            if (updatedTopic) {
                setSelectedTopic({...updatedTopic, subscribed: true})
            }
        }
    }

    // Handle topic selection
    const handleTopicSelect = (topic: Topic) => {
        setSelectedTopic(topic)
    }

    const handleToggleConnect = () => {
        setIsConnected((prev) => !prev);
    };
    // Handle topic unsubscription
    const handleUnsubscribe = (topicId: string) => {
        unsubscribeTopic(topicId)
        // Update the selected topic if it's the one being unsubscribed from
        if (selectedTopic && selectedTopic.id === topicId) {
            const updatedTopic = topics.find((t) => t.id === topicId)
            if (updatedTopic) {
                setSelectedTopic({...updatedTopic, subscribed: false})
            }
        }
    }

    // Handle search input
    const handleSearch = (query: string) => {
        setSearchQuery(query)
    }

    return (
        <div className={"flex flex-col h-screen"}>
            <Header onSearch={handleSearch}/>
            <div className={"flex flex-1 overflow-hidden"}>
                <TopicPanel
                    topics={filteredTopics}
                    selectedTopic={selectedTopic}
                    onSelectTopic={handleTopicSelect}
                    onCreateTopic={createTopic}
                    onDeleteTopic={deleteTopic}
                    onRenameTopic={renameTopic}
                    onSubscribe={handleSubscribe}
                    onUnsubscribe={handleUnsubscribe}
                    isConnected={isConnected}
                    onToggleConnect={handleToggleConnect}
                />
                {isConnected ? (
                    selectedTopic ? (
                        <MessagePanel
                            topic={selectedTopic}
                            messages={messages.filter((m) => m.topic === selectedTopic.name)}
                            onPublish={publishMessage}
                            onSubscribe={handleSubscribe}
                            onUnsubscribe={handleUnsubscribe}
                        />
                    ) : (
                        <div className="flex-1 flex items-center justify-center border-t">
                            <p>Select a topic to view messages</p>
                        </div>
                    )
                ) : (
                    <ConnectionPanel onToggleConnect={handleToggleConnect}/>
                )}

            </div>
        </div>
    )
}

export default MqttDashboard
