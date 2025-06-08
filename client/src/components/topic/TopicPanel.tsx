import React, {useEffect, useState} from "react";
import CreateTopicForm from "./CreateTopicForm";
import TopicItem from "./TopicItem";
import {Button} from "@/components/ui/button.tsx";
import {MessageSquare, Plus} from "lucide-react";
import {ScrollArea} from "../ui/scroll-area";
import type {Topic} from "@/lib/types.ts";
import {useTopics} from "@/api/hooks/useTopics";
import {useSubscribeTopics} from "@/api/hooks/useSubscribeTopic.ts";

interface TopicPanelProps {
    selectedTopic: Topic | null;
    onSelectTopic: (topic: any) => void;
    onSubscribe: (id: string) => void;
    onUnsubscribe: (id: string) => void;
    isConnected: boolean;
    searchTerm: string;
}

export default function TopicPanel({
                                       selectedTopic,
                                       onSelectTopic,
                                       isConnected,
                                       searchTerm,
                                   }: TopicPanelProps) {
    const TOPIC_CHECK_DELAY = 200;
    const [isCreating, setIsCreating] = useState(false);
    const [newTopicName, setNewTopicName] = useState("");
    const [editingTopicId, setEditingTopicId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");

    const {fetchTopics, topics, addTopic, error, isLoading} = useTopics();
    const {
        subscribeToTopics,
        isLoading: isSubscribing,
        error: subscribeError,
        clearError: clearSubscribeError,
    } = useSubscribeTopics();

    useEffect(() => {
        fetchTopics().catch((error) => {
            console.error("Failed to fetch topics.", error);
        });
        console.log("Fetched Topics: ", topics);
    }, [fetchTopics]);

    const handleCreateSubmit = async (e: React.FormEvent): Promise<boolean> => {
        e.preventDefault();
        const trimmed = newTopicName.trim();
        if (!trimmed) return false;
        clearSubscribeError();
        try {
            const result = await subscribeToTopics([trimmed]);
            const topicResult = result.results.find(r => r.topic === trimmed);

            if (topicResult?.success) {
                return handleSuccessfulSubscription(trimmed);
            }
        } catch (error) {
            console.error("Subscription failed:", error);
        }
        return await checkTopicCreationWithFallback(trimmed);
    };

    const handleSuccessfulSubscription = (topicName: string): boolean => {
        console.log("Topic created and subscribed successfully!");

        const brokerId = localStorage.getItem('brokerId') || '';
        const userId = localStorage.getItem('userId') || '';
        const newTopic: Topic = {
            Id: topicName,
            Topic: topicName,
            BrokerId: brokerId,
            UserId: userId,
            CreationDate: new Date().toISOString(),
        };
        addTopic(newTopic);
        setNewTopicName("");
        fetchTopics().catch(console.error);
        return true;
    };

    const checkTopicCreationWithFallback = async (
        topicName: string
    ): Promise<boolean> => {
        try {
            await fetchTopics();
            return new Promise((resolve) => {
                setTimeout(() => {
                    const wasCreated = topics.some(topic =>
                        topic.Topic.toLowerCase() === topicName.toLowerCase()
                    );
                    if (wasCreated) {
                        setNewTopicName("");
                    }
                    resolve(wasCreated);
                }, TOPIC_CHECK_DELAY);
            });
        } catch (error) {
            console.error("Failed to refresh topics:", error);
            return false;
        }
    };

    const handleClearTopic = () => {
        onSelectTopic(null);
        const url = new URL(window.location.href);
        url.searchParams.delete("topic");
        window.history.pushState({}, "", url.toString());
    };

    const handleTopicSelect = (topic: Topic) => {
        // If clicking the already selected topic, deselect it
        if (selectedTopic?.Id === topic.Id) {
            handleClearTopic();
        } else {
            onSelectTopic(topic);
            const url = new URL(window.location.href);
            url.searchParams.set("topic", topic.Topic);
            window.history.pushState({}, "", url.toString());
        }
    };

    const startEditing = (topic: Topic) => {
        setEditingTopicId(topic.Id);
        setEditingName(topic.Topic);
    };

    const handleRetry = () => {
        fetchTopics().catch((error) => {
            console.error("Failed to refetch topics.", error);
        });
        console.log("Refetched Topics: ", topics);
    };

    const handleCreateCancel = () => {
        setIsCreating(false);
        setNewTopicName("");
        clearSubscribeError();
    };

    const filteredTopics = React.useMemo(() => {
        if (!searchTerm.trim()) {
            return topics;
        }
        return topics.filter(topic =>
            topic.Topic.toLowerCase().includes(searchTerm.toLowerCase())
        );
    }, [topics, searchTerm]);

    const renderTopicsContent = () => {
        if (isLoading) {
            return (
                <div className="p-2">
                    <div className="text-center py-8">
                        <div
                            className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100 mx-auto mb-2"></div>
                        <p className="text-sm text-gray-500">Loading topics...</p>
                    </div>
                </div>
            );
        }

        if (error) {
            return (
                <div className="p-2">
                    <div className="text-center py-8">
                        <div className="text-red-500 mb-2">
                            <svg
                                className="h-12 w-12 mx-auto mb-2"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"
                                />
                            </svg>
                        </div>
                        <p className="text-red-600 dark:text-red-400 mb-2 font-medium">
                            Failed to load topics
                        </p>
                        <p className="text-sm text-gray-500 mb-4">
                            {error.message || "An unexpected error occurred"}
                        </p>
                        <Button
                            size="sm"
                            variant="outline"
                            onClick={handleRetry}
                            disabled={isLoading}
                        >
                            Try Again
                        </Button>
                    </div>
                </div>
            );
        }

        if (!Array.isArray(filteredTopics) || filteredTopics.length === 0) {
            return (
                <div className="p-2">
                    <div className="text-center py-8">
                        <MessageSquare className="h-12 w-12 mx-auto mb-2 opacity-20"/>
                        <p className="text-sm text-gray-500">
                            No topics found. Create a new topic to get started.
                        </p>
                    </div>
                </div>
            );
        }

        return (
            <div className="p-2">
                <ul className="space-y-1">
                    {filteredTopics.map((topic) => (
                        <TopicItem
                            key={topic.Id}
                            topic={topic}
                            isEditing={editingTopicId === topic.Id}
                            editingName={editingName}
                            selected={selectedTopic?.Id === topic.Id}
                            onSelect={() => handleTopicSelect(topic)}
                            onStartEditing={startEditing}
                            onEditNameChange={(e) => setEditingName(e.target.value)}
                        />
                    ))}
                </ul>
            </div>
        );
    };

    return (
        <div className="w-100 border-r border-t h-full flex flex-col">
            <div className="p-4 h-16 border-b flex items-center justify-between flex-shrink-0"></div>
            {isConnected && (
                <div
                    className="p-4 h-16 flex items-center justify-between bg-[var(--background)] border-b border-[var(--border)] flex-shrink-0">
                    <h2 className="dark:text-gray-200 flex items-center">
                        Topics
                        {isLoading && (
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-500 ml-2"></div>
                        )}
                    </h2>
                    <Button
                        size="sm"
                        onClick={() => {
                            setIsCreating(true);
                            setNewTopicName("");
                            clearSubscribeError();
                        }}
                        disabled={isLoading || !!error}
                    >
                        <Plus className="h-4 w-4 mr-1"/>
                        New Topic
                    </Button>
                </div>
            )}
            {isConnected && (
                <div className="flex-1 flex flex-col overflow-hidden">
                    {isCreating && !isLoading && !error && (
                        <div className="flex-shrink-0">
                            <CreateTopicForm
                                newTopicName={newTopicName}
                                onNewTopicNameChange={(e) => setNewTopicName(e.target.value)}
                                onSubmit={handleCreateSubmit}
                                onCancel={handleCreateCancel}
                                isSubmitting={isSubscribing}
                                error={subscribeError}
                                onClearError={clearSubscribeError}
                            />
                        </div>
                    )}
                    <div className="flex-1 overflow-hidden">
                        <ScrollArea className="h-full">
                            {renderTopicsContent()}
                        </ScrollArea>
                    </div>
                </div>
            )}
        </div>
    );
}