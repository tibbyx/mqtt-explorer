import React, {useEffect, useState} from "react";
import CreateTopicForm from "./CreateTopicForm";
import TopicItem from "./TopicItem";
import {Button} from "@/components/ui/button.tsx";
import {MessageSquare, Plus} from "lucide-react";
import {ScrollArea} from "../ui/scroll-area";
import type {Topic} from "@/lib/types.ts";

interface TopicPanelProps {
    topics: Topic[];
    selectedTopic: Topic | null;
    onSelectTopic: (topic: Topic) => void;
    onCreateTopic: (name: string) => void;
    onDeleteTopic: (id: string) => void;
    onRenameTopic: (id: string, newName: string) => void;
    onSubscribe: (id: string) => void;
    onUnsubscribe: (id: string) => void;
}

export default function TopicPanel({
                                       topics,
                                       selectedTopic,
                                       onSelectTopic,
                                       onCreateTopic,
                                       onDeleteTopic,
                                       onRenameTopic,
                                   }: TopicPanelProps) {
    const [isCreating, setIsCreating] = useState(false);
    const [newTopicName, setNewTopicName] = useState("");
    const [editingTopicId, setEditingTopicId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");
    const [localTopics, setLocalTopics] = useState<Topic[]>([]);
    const [showTopics, setShowTopics] = useState(false);
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        setLocalTopics(topics);
    }, [topics]);

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
            className="w-100 border-r border-gray-200 dark:border-gray-800 flex flex-col h-full bg-white dark:bg-gray-950">
            <div className="p-4 border-b border-gray-200 dark:border-gray-800 flex items-center justify-between">
                <h2 className="font-semibold text-gray-700 dark:text-gray-200">
                    Server
                </h2>
                <Button
                    size="sm"
                    onClick={() => {
                        if (connected) {
                            // Disconnect handler hier
                            console.log("Disconnect clicked");
                            setConnected(false);
                        } else {
                            // Connect handler hier
                            console.log("Connect clicked");
                            setConnected(true);
                        }
                    }}
                    className="bg-[#7a62f6] hover:bg-[#6952e3] text-white rounded-full"
                >
                    {connected ? "Disconnect" : "Connect"}
                </Button>
            </div>

            <div
                className="p-4 border-b border-gray-200 dark:border-gray-800 flex items-center justify-between cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-900 transition"
                onClick={() => setShowTopics(prev => !prev)}
            >
                <h2 className="font-semibold text-gray-700 dark:text-gray-200">
                    Topics
                </h2>
                <Button
                    size="sm"
                    onClick={(e) => {
                        e.stopPropagation();
                        setIsCreating(true);
                        setNewTopicName("");
                    }}
                    className="bg-[#7a62f6] hover:bg-[#6952e3] text-white rounded-full"
                >
                    <Plus className="h-4 w-4 mr-1"/>
                    New Topic
                </Button>
            </div>


            {showTopics && (
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
                                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                                    <MessageSquare className="h-12 w-12 mx-auto mb-2 opacity-20"/>
                                    <p>No topics found. Create a new topic to get started.</p>
                                </div>
                            ) : (
                                <ul className="space-y-1">
                                    {localTopics.map((topic) => (
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
