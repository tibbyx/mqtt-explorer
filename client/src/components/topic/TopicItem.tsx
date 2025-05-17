import React from "react";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "../ui/button";
import {Check, Edit, Trash2, X} from "lucide-react";
import type {Topic} from "@/lib/types.ts";

interface TopicItemProps {
    topic: Topic;
    isEditing: boolean;
    editingName: string;
    selected: boolean;
    onSelect: () => void;
    onStartEditing: (topic: Topic) => void;
    onDelete: () => void;
    onEditNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSubmitEdit: () => void;
    onCancelEdit: () => void;
}

export default function TopicItem({
                                      topic,
                                      isEditing,
                                      editingName,
                                      selected,
                                      onSelect,
                                      onStartEditing,
                                      onDelete,
                                      onEditNameChange,
                                      onSubmitEdit,
                                      onCancelEdit,
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
                    <Button
                        size="sm"
                        className="h-8 w-8 p-0"
                        onClick={onSubmitEdit}
                    >
                        <Check className="h-4 w-4"/>
                    </Button>
                    <Button
                        size="sm"
                        variant="outline"
                        className="h-8 w-8 p-0"
                        onClick={onCancelEdit}
                    >
                        <X className="h-4 w-4"/>
                    </Button>
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
                <div className="flex items-center gap-2 flex-1 min-w-0">
          <span
              className={`truncate ${
                  selected ? "font-medium" : ""
              }`}
          >
            {topic.name}
          </span>
                    {topic.subscribed && (
                        <span
                            className="text-xs px-2 py-0.5 rounded-full font-medium bg-[var(--background)] border border-[var(--border)]">
              Subscribed
            </span>
                    )}
                </div>
                <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 w-7 p-0 hover:bg-[#7a62f6]/10"
                        onClick={(e) => {
                            e.stopPropagation();
                            onStartEditing(topic);
                        }}
                    >
                        <Edit className="h-3.5 w-3.5"/>
                    </Button>
                    <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 w-7 p-0 text-gray-500 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg"
                        onClick={(e) => {
                            e.stopPropagation();
                            onDelete();
                        }}
                    >
                        <Trash2 className="h-3.5 w-3.5"/>
                    </Button>
                </div>
            </div>
        </li>
    );
}
