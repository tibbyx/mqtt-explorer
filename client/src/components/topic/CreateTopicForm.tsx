import React from "react";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Check, X} from "lucide-react";

interface CreateTopicFormProps {
    newTopicName: string;
    onNewTopicNameChange: (
        e: React.ChangeEvent<HTMLInputElement>
    ) => void;
    onSubmit: (e: React.FormEvent) => void;
    onCancel: () => void;
}

export default function CreateTopicForm({
                                            newTopicName,
                                            onNewTopicNameChange,
                                            onSubmit,
                                            onCancel,
                                        }: CreateTopicFormProps) {
    return (
        <form
            onSubmit={onSubmit}
            className="p-3 border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-900"
        >
            <div className="flex gap-2">
                <Input
                    placeholder="Topic name"
                    value={newTopicName}
                    onChange={onNewTopicNameChange}
                    autoFocus
                    className="h-9 rounded-lg border-gray-200 dark:border-gray-800 focus-visible:ring-[#7a62f6]"
                />
                <Button
                    type="submit"
                    size="sm"
                    className="h-9 px-2 bg-[#7a62f6] hover:bg-[#6952e3] rounded-lg"
                >
                    <Check className="h-4 w-4"/>
                </Button>
                <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    className="h-9 px-2 rounded-lg border-gray-200 dark:border-gray-800"
                    onClick={onCancel}
                >
                    <X className="h-4 w-4"/>
                </Button>
            </div>
        </form>
    );
}
