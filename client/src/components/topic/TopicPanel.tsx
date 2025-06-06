import React, {useEffect, useState} from "react";
import CreateTopicForm from "./CreateTopicForm";
import TopicItem from "./TopicItem";
import {Button} from "@/components/ui/button.tsx";
import {MessageSquare, Plus} from "lucide-react";
import {ScrollArea} from "../ui/scroll-area";
import type {Topic} from "@/lib/types.ts";
import {ChevronDown} from "lucide-react";
import {useTopics} from "@/api/hooks/useTopics";

interface TopicPanelProps {
    topics: Topic[];
    selectedTopic: Topic | null;
    onSelectTopic: (topic: Topic) => void;
    onCreateTopic: (name: string) => void;
    onDeleteTopic: (id: string) => void;
    onRenameTopic: (id: string, newName: string) => void;
    onSubscribe: (id: string) => void;
    onUnsubscribe: (id: string) => void;
    onToggleConnect: () => void;
    isConnected: boolean;
}

export default function TopicPanel({
                                       /*topics*/
                                       selectedTopic,
                                       onSelectTopic,
                                       onCreateTopic,
                                       onDeleteTopic,
                                       onRenameTopic,
                                       onToggleConnect,
                                       isConnected,
                                   }: TopicPanelProps) {
    const [isCreating, setIsCreating] = useState(false);
    const [newTopicName, setNewTopicName] = useState("");
    const [editingTopicId, setEditingTopicId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");
    const [localTopics, setLocalTopics] = useState<Topic[]>([]);
    const [showTopics, setShowTopics] = useState(true);
    const {fetchTopics, topics, error, isLoading} = useTopics()

    /*    useEffect(() => {
            setLocalTopics(topics);
        }, [topics]);*/

    useEffect(() => {
        fetchTopics()
        console.log("Fetched Topics: ", topics)
    }, [fetchTopics]);

    const handleCreateSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const trimmed = newTopicName.trim();
        if (trimmed) {
            onCreateTopic(trimmed);
            setNewTopicName("");
            setIsCreating(false);
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

    return (
        <div
            className="w-100 border-r border-t h-full">
            <div className="p-4 h-16 border-b flex items-center justify-between">
                <h2>
                    Server
                </h2>
                {isConnected && (
                    <Button
                        size="sm"
                        onClick={onToggleConnect}
                    >
                        Disconnect
                    </Button>
                )}
            </div>
            {isConnected && (
                <div
                    className="p-4 h-16 flex items-center justify-between cursor-pointer transition bg-[var(--background)] border-b border-[var(--border)]"
                    onClick={() => setShowTopics(prev => !prev)}
                >
                    <h2 className="dark:text-gray-200 ">
                        Topics
                    </h2>
                    {!showTopics && (
                        <ChevronDown className="w-4 h-4 ml-1"/>
                    )}
                    {showTopics && (
                        <Button
                            size="sm"
                            onClick={(e) => {
                                e.stopPropagation();
                                setIsCreating(true);
                                setNewTopicName("");
                            }}
                        >
                            <Plus className="h-4 w-4 mr-1"/>
                            New Topic
                        </Button>
                    )}
                </div>
            )}
            {isConnected && showTopics && (
                <>
                    {isCreating && (
                        <CreateTopicForm
                            newTopicName={newTopicName}
                            onNewTopicNameChange={(e) => setNewTopicName(e.target.value)}
                            onSubmit={handleCreateSubmit}
                            onCancel={() => setIsCreating(false)}
                        />
                    )}
                    <ScrollArea className="flex-1">
                        <div className="p-2">
                            {localTopics.length === 0 ? (
                                <div className="text-center py-8">
                                    <MessageSquare className="h-12 w-12 mx-auto mb-2 opacity-20"/>
                                    <p>No topics found. Create a new topic to get started.</p>
                                </div>
                            ) : (
                                <ul className="space-y-1">
                                    {topics.map((topic) => (
                                        <TopicItem
                                            key={topic.id}
                                            topic={topic}
                                            isEditing={editingTopicId === topic.id}
                                            editingName={editingName}
                                            selected={selectedTopic?.id === topic.id}
                                            onSelect={() => onSelectTopic(topic)}
                                            onStartEditing={startEditing}
                                            onDelete={() => onDeleteTopic(topic.id)}
                                            onEditNameChange={(e) => setEditingName(e.target.value)}
                                            onSubmitEdit={() => handleEditSubmit(topic.id)}
                                            onCancelEdit={cancelEditing}
                                        />
                                    ))}
                                </ul>
                            )}
                        </div>
                    </ScrollArea>
                </>
            )}
        </div>
    );
}
