import React from "react";
import {Input} from "@/components/ui/input.tsx";
import type {Topic} from "@/lib/types.ts";

interface TopicItemProps {
    topic: Topic;
    isEditing: boolean;
    editingName: string;
    selected: boolean;
    onSelect: () => void;
    onStartEditing: (topic: Topic) => void;
    onEditNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export default function TopicItem({
                                      topic,
                                      isEditing,
                                      editingName,
                                      selected,
                                      onSelect,
                                      onEditNameChange,
                                  }: TopicItemProps) {
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
                {1 === 1 && (
                    <span
                        className="text-xs px-2 py-0.5 rounded-full font-medium bg-[var(--background)] border border-[var(--border)] ml-2 flex-shrink-0">
                    Subscribed
                </span>
                )}
            </div>
        </li>
    );
}
