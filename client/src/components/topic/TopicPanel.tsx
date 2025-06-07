import React, {useEffect, useState} from "react";
import CreateTopicForm from "./CreateTopicForm";
import TopicItem from "./TopicItem";
import {Button} from "@/components/ui/button.tsx";
import {MessageSquare, Plus} from "lucide-react";
import {ScrollArea} from "../ui/scroll-area";
import type {Topic} from "@/lib/types.ts";
import {ChevronDown} from "lucide-react";
import {useTopics} from "@/api/hooks/useTopics";
import {useSubscribeTopics} from "@/api/hooks/useSubscribeTopic.ts";

interface TopicPanelProps {
    topics: Topic[];
    selectedTopic: Topic | null;
    onSelectTopic: (topic: any) => void;
    onCreateTopic: (name: string) => void;
    onDeleteTopic: (id: string) => void;
    onRenameTopic: (id: string, newName: string) => void;
    onSubscribe: (id: string) => void;
    onUnsubscribe: (id: string) => void;
    isConnected: boolean;
}

export default function TopicPanel({
                                       selectedTopic,
                                       onSelectTopic,
                                       onCreateTopic,
                                       onDeleteTopic,
                                       onRenameTopic,
                                       isConnected,
                                   }: TopicPanelProps) {
    const [isCreating, setIsCreating] = useState(false);
    const [newTopicName, setNewTopicName] = useState("");
    const [editingTopicId, setEditingTopicId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");
    const [showTopics, setShowTopics] = useState(true);

    const {fetchTopics, topics, addTopic, error, isLoading} = useTopics();
    const {
        subscribeToTopics,
        isLoading: isSubscribing,
    } = useSubscribeTopics();

    useEffect(() => {
        fetchTopics().catch((error) => {
            console.error("Failed to fetch topics.", error);
        });
        console.log("Fetched Topics: ", topics);
    }, [fetchTopics]);

    const handleCreateSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const trimmed = newTopicName.trim();
        if (!trimmed) return;

        console.log("Before submit - isCreating:", isCreating); // Debug

        try {
            const response = await subscribeToTopics([trimmed]);
            const result = response.result[trimmed];

            if (result.status === "Fine") {
                const newTopic: Topic = {
                    id: trimmed,
                    name: trimmed,
                    subscribed: true,
                };

                addTopic(newTopic);
                onCreateTopic(trimmed);

                console.log("Setting isCreating to false"); // Debug
                setNewTopicName("");
                setIsCreating(false);

                console.log("After submit - isCreating should be false"); // Debug
            } else {
                console.error("Failed to create topic:", result.message);
            }
        } catch (error) {
            console.error("Error creating topic:", error);
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
        if (selectedTopic?.id === topic.id) {
            handleClearTopic();
        } else {
            onSelectTopic(topic);
            const url = new URL(window.location.href);
            url.searchParams.set("topic", topic.name);
            window.history.pushState({}, "", url.toString());
        }
    };

    const handleEditSubmit = (id: string) => {
        const trimmed = editingName.trim();
        if (trimmed) {
            onRenameTopic(id, trimmed);
        }
        setEditingTopicId(null);
        setEditingName("");
    };

    const startEditing = (topic: Topic) => {
        setEditingTopicId(topic.id);
        setEditingName(topic.name);
    };

    const cancelEditing = () => {
        setEditingTopicId(null);
        setEditingName("");
    };

    const handleRetry = () => {
        fetchTopics().catch((error) => {
            console.error("Failed to refetch topics.", error);
        });
        console.log("Refetched Topics: ", topics);
    };

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

        if (topics.length === 0) {
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
                    {topics.map((topic) => (
                        <TopicItem
                            key={topic.id}
                            topic={topic}
                            isEditing={editingTopicId === topic.id}
                            editingName={editingName}
                            selected={selectedTopic?.id === topic.id}
                            onSelect={() => handleTopicSelect(topic)}
                            onStartEditing={startEditing}
                            onDelete={() => onDeleteTopic(topic.id)}
                            onEditNameChange={(e) => setEditingName(e.target.value)}
                            onSubmitEdit={() => handleEditSubmit(topic.id)}
                            onCancelEdit={cancelEditing}
                        />
                    ))}
                </ul>
            </div>
        );
    };

    return (
        <div className="w-100 border-r border-t h-full">
            <div className="p-4 h-16 border-b flex items-center justify-between"></div>
            {isConnected && (
                <div
                    className="p-4 h-16 flex items-center justify-between cursor-pointer transition bg-[var(--background)] border-b border-[var(--border)]"
                    onClick={() => setShowTopics((prev) => !prev)}
                >
                    <h2 className="dark:text-gray-200 flex items-center">
                        Topics
                        {isLoading && (
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-500 ml-2"></div>
                        )}
                    </h2>
                    {!showTopics && <ChevronDown className="w-4 h-4 ml-1"/>}
                    {showTopics && (
                        <Button
                            size="sm"
                            onClick={(e) => {
                                e.stopPropagation();
                                setIsCreating(true);
                                setNewTopicName("");
                            }}
                            disabled={isLoading || !!error}
                        >
                            <Plus className="h-4 w-4 mr-1"/>
                            New Topic
                        </Button>
                    )}
                </div>
            )}
            {isConnected && showTopics && (
                <>
                    {isCreating && !isLoading && !error && (
                        <CreateTopicForm
                            newTopicName={newTopicName}
                            onNewTopicNameChange={(e) => setNewTopicName(e.target.value)}
                            onSubmit={handleCreateSubmit}
                            onCancel={() => {
                                setIsCreating(false);
                                setNewTopicName("");
                            }}
                            isSubmitting={isSubscribing}
                        />
                    )}
                    <ScrollArea className="flex-1">{renderTopicsContent()}</ScrollArea>
                </>
            )}
        </div>
    );
}