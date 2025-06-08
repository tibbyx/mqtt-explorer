import {Header} from "@/components/Header.tsx";
import TopicPanel from "@/components/topic/TopicPanel.tsx";
import {MessagePanel} from "@/components/message/MessagePanel.tsx";
import {ConnectionPanel} from "@/components/connection/ConnectionPanel.tsx";
import type {Topic} from "@/lib/types.ts";
import {useEffect, useState} from "react";
import {useMessages} from "@/api/hooks/useMessages.ts";
import {useTopics} from "@/api/hooks/useTopics.ts";

function MqttDashboard() {
    const [selectedTopic, setSelectedTopic] = useState<Topic | null>(null)
    const [isConnected, setIsConnected] = useState(false);
    const {topics} = useTopics();
    const {messages, fetchMessages, clearMessages} = useMessages();
    const [searchTerm, setSearchTerm] = useState("");

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const topicFromUrl = urlParams.get('topic');

        if (topicFromUrl && topics.length > 0 && isConnected) {
            const topic = topics.find(t => t.Topic === topicFromUrl);
            if (topic) {
                setSelectedTopic(topic);
                // Load messages for this topic
                fetchMessages(topic.Topic).catch(console.error);
            }
        }
    }, [topics, isConnected, fetchMessages]);

    useEffect(() => {
        if (selectedTopic && isConnected) {
            fetchMessages(selectedTopic.Topic).catch(console.error);
        } else {
            clearMessages();
        }
    }, [selectedTopic, isConnected, fetchMessages, clearMessages]);

    const handleTopicSelect = (topic: Topic | null) => {
        setSelectedTopic(topic);

        if (topic) {
            const url = new URL(window.location.href);
            url.searchParams.set('topic', topic.Topic);
            window.history.pushState({}, '', url.toString());

            fetchMessages(topic.Topic).catch(console.error);
        } else {
            const url = new URL(window.location.href);
            url.searchParams.delete('topic');
            window.history.pushState({}, '', url.toString());
            clearMessages();
        }
    };

    const handleSubscribe = (topicId: string) => {
        if (selectedTopic && selectedTopic.Id === topicId) {
            const updatedTopic = topics.find((t: Topic) => t.Id === topicId)
            if (updatedTopic) {
                setSelectedTopic({...updatedTopic})
            }
        }
    }

    const handleToggleConnect = () => {
        setIsConnected((prev) => !prev);
    };
    const handleUnsubscribe = (topicId: string) => {
        if (selectedTopic && selectedTopic.Id === topicId) {
            const updatedTopic = topics.find((t) => t.Id === topicId)
            if (updatedTopic) {
                setSelectedTopic({...updatedTopic})
            }
        }
    }

    const handleSearch = (searchTerm: string) => {
        setSearchTerm(searchTerm);
        console.log("Searching for:", searchTerm);
    };

    return (
        <div className={"flex flex-col h-screen"}>
            <Header
                isConnected={isConnected}
                onToggleConnect={handleToggleConnect}
                onSearch={handleSearch}
            />
            <div className={"flex flex-1 overflow-hidden"}>
                {isConnected ? (
                    <>
                        <TopicPanel
                            selectedTopic={selectedTopic}
                            onSelectTopic={handleTopicSelect}
                            onSubscribe={handleSubscribe}
                            onUnsubscribe={handleUnsubscribe}
                            isConnected={isConnected}
                            searchTerm={searchTerm}
                        />
                        {selectedTopic ? (
                            <MessagePanel
                                topic={selectedTopic}
                                messages={messages.filter((m) => m.topic === selectedTopic.Topic)}
                                onSubscribe={handleSubscribe}
                                onUnsubscribe={handleUnsubscribe}
                            />
                        ) : (
                            <div className="flex-1 flex items-center justify-center border-t">
                                <p>Select a topic to view messages</p>
                            </div>
                        )}
                    </>
                ) : (
                    <ConnectionPanel onToggleConnect={handleToggleConnect}/>
                )}
            </div>
        </div>
    );
}

export default MqttDashboard
