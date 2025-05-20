import { Header } from "@/components/Header.tsx";
import TopicPanel from "@/components/topic/TopicPanel.tsx";
import { MessagePanel } from "@/components/message/MessagePanel.tsx";
import { ConnectionPanel } from "@/components/connection/ConnectionPanel.tsx";
import type { Topic } from "@/lib/types";
import { useMqttWebSocket } from "@/hooks/use-mqtt-websocket.tsx";
import { useToast } from "../hooks/use-toast";
import { useState } from "react";

function MqttDashboard() {
    const [selectedTopics, setSelectedTopics] = useState<Topic[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [isConnected, setIsConnected] = useState(false);
    const { toast } = useToast();
    const [isSplitScreen, setIsSplitScreen] = useState(false);

    const {
        messages,
        createTopic,
        deleteTopic,
        renameTopic,
        subscribeTopic,
        unsubscribeTopic,
        filteredTopics,
        publishMessage,
    } = useMqttWebSocket({
        shouldConnect: isConnected,
        onConnect: () => {
            setIsConnected(true);
            toast({
                title: "Connected",
                description: "Successfully connected to MQTT broker",
                className: "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-900/50",
            });
        },
        onDisconnect: () => {
            setIsConnected(false);
            toast({
                title: "Disconnected",
                description: "Connection to MQTT broker lost",
                variant: "destructive",
            });
            setIsSplitScreen(false);
            setSelectedTopics([]);
        },
        searchQuery,
    });

    const handleSubscribe = (topicId: string) => {
        subscribeTopic(topicId);
        setSelectedTopics(prev =>
            prev.map(t => (t.id === topicId ? { ...t, subscribed: true } : t))
        );
    };

    const handleUnsubscribe = (topicId: string) => {
        unsubscribeTopic(topicId);
        setSelectedTopics(prev =>
            prev.map(t => (t.id === topicId ? { ...t, subscribed: false } : t))
        );
    };

    const handleTopicSelect = (topic: Topic) => {
        console.log('Angeklicktes Topic ' + topic.name);
        setSelectedTopics(prev => {
            if (!isSplitScreen) {
                return [topic];
            } else {

                const [mainTopic, secondaryTopic] = prev;

                if (mainTopic?.id === topic.id) return prev;

                if (secondaryTopic?.id === topic.id) return prev;

                return [mainTopic, topic];
            }
        });
    };

    const handleCloseTopic = (topicId: string) => {
        setSelectedTopics(prev => {
            const filtered = prev.filter(t => t.id !== topicId);
            if (filtered.length < 2) setIsSplitScreen(false);
            return filtered;
        });
    };
    
    const handleToggleSplitScreen = () => {
        setIsSplitScreen(prev => {
            const newSplit = !prev;

            if (newSplit && selectedTopics.length === 1) {
                setSelectedTopics(prevTopics => [...prevTopics, null].filter(Boolean) as Topic[]);
            }

            if (!newSplit) {
                setSelectedTopics(prevTopics => (prevTopics.length > 0 ? [prevTopics[0]] : []));
            }

            return newSplit;
        });
    };

    const handleToggleConnect = () => {
        setIsConnected(prev => !prev);
    };

    const handleSearch = (query: string) => {
        setSearchQuery(query);
    };

    return (
        <div className="flex flex-col h-screen">
            <Header onSearch={handleSearch} isConnected={isConnected}/>
            <div className="flex flex-1 overflow-hidden">
                <TopicPanel
                    topics={filteredTopics}
                    selectedTopic={selectedTopics[0] || null}
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
                    selectedTopics.length === 0 ? (
                        <div className="flex-1 flex items-center justify-center border-t">
                            <p>Select a topic to view messages</p>
                        </div>
                    ) : isSplitScreen && selectedTopics.length === 2 ? (
                        <div className="flex flex-row flex-1 h-full">
                            {/* Left panel */}
                            <MessagePanel
                                key={selectedTopics[0].id}
                                topic={selectedTopics[0]}
                                messages={messages.filter(m => m.topic === selectedTopics[0].name)}
                                onPublish={publishMessage}
                                onSubscribe={handleSubscribe}
                                onUnsubscribe={handleUnsubscribe}
                                onCloseTopic={() => handleCloseTopic(selectedTopics[0].id)}
                                isSplitScreen={isSplitScreen}
                                onToggleSplitScreen={handleToggleSplitScreen}
                                className="w-1/2"
                            />

                            {/* Vertical Divider */}
                            <div className="w-px bg-gray-300 dark:bg-gray-600"/>

                            {/* Right panel */}
                            <MessagePanel
                                key={selectedTopics[1].id}
                                topic={selectedTopics[1]}
                                messages={messages.filter(m => m.topic === selectedTopics[1].name)}
                                onPublish={publishMessage}
                                onSubscribe={handleSubscribe}
                                onUnsubscribe={handleUnsubscribe}
                                onCloseTopic={() => handleCloseTopic(selectedTopics[1].id)}
                                isSplitScreen={isSplitScreen}
                                onToggleSplitScreen={handleToggleSplitScreen}
                                className="w-1/2"
                            />
                        </div>
                    ) : (
                        <MessagePanel
                            topic={selectedTopics[0]}
                            messages={messages.filter(m => m.topic === selectedTopics[0].name)}
                            onPublish={publishMessage}
                            onSubscribe={handleSubscribe}
                            onUnsubscribe={handleUnsubscribe}
                            onCloseTopic={() => handleCloseTopic(selectedTopics[0].id)}
                            isSplitScreen={isSplitScreen}
                            onToggleSplitScreen={handleToggleSplitScreen}
                        />
                    )
                ) : (
                    <ConnectionPanel onToggleConnect={handleToggleConnect}/>
                )}
            </div>
        </div>
    );
}

export default MqttDashboard;
