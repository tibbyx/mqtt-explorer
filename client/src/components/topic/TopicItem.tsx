import React, {useState, useEffect} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import type {Topic} from "@/lib/types.ts";
import {useTopicSubscription} from "@/api/hooks/useTopicSubscription.ts";

interface TopicItemProps {
    topic: Topic;
    isEditing: boolean;
    editingName: string;
    selected: boolean;
    onSelect: () => void;
    onStartEditing: (topic: Topic) => void;
    onEditNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSubscriptionChange?: (topicId: number, subscribed: boolean) => void;
}

export default function TopicItem({
                                      topic,
                                      isEditing,
                                      editingName,
                                      selected,
                                      onSelect,
                                      onEditNameChange,
                                      onSubscriptionChange,
                                  }: TopicItemProps) {
    const {isLoading, toggleSubscription} = useTopicSubscription();
    const [localSubscribed, setLocalSubscribed] = useState(topic.Subscribed);

    useEffect(() => {
        setLocalSubscribed(topic.Subscribed);
    }, [topic.Subscribed]);

    const handleToggleSubscription = async (e: React.MouseEvent) => {
        e.stopPropagation();

        const newState = !localSubscribed;
        setLocalSubscribed(newState);

        try {
            const actualNewState = await toggleSubscription(
                topic.Topic,
                localSubscribed
            );

            setLocalSubscribed(actualNewState);
            onSubscriptionChange?.(Number(topic.Id), actualNewState);
        } catch (error) {
            console.error('Failed to toggle subscription:', error);
            setLocalSubscribed(topic.Subscribed);
        }
    };

    if (isEditing) {
        return (
            <li>
                <div className="flex items-center gap-1 p-1">
                    <Input
                        value={editingName}
                        onChange={onEditNameChange}
                        autoFocus
                        className="h-8 flex-1"
                    />
                </div>
            </li>
        );
    }

    return (
        <li>
            <div
                className={`flex items-center justify-between p-2.5 rounded-xl transition-all duration-200 cursor-pointer group ${
                    selected
                        ? "bg-[var(--primary)]/20 dark:bg-[var(--primary)]/90"
                        : "hover:bg-[var(--primary)]/10 dark:hover:bg-[var(--primary)]/30"
                }`}
                onClick={onSelect}
            >
                <div className="flex items-center flex-1 min-w-0">
                    <span
                        className={`truncate ${
                            selected ? "font-medium" : ""
                        }`}
                    >
                    {topic.Topic}
                </span>
                </div>
                <Button
                    size="sm"
                    variant="ghost"
                    onClick={handleToggleSubscription}
                    disabled={isLoading}
                    className="text-xs px-2 py-0.5 h-auto rounded-full font-medium bg-[var(--background)] border border-[var(--border)] flex-shrink-0 hover:bg-[var(--background)]/80"
                >
                    {isLoading ? "..." : localSubscribed ? "Subscribed" : "Unsubscribed"}
                </Button>
            </div>
        </li>
    );
}