import React from "react";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Check, X, AlertCircle} from "lucide-react";

interface CreateTopicFormProps {
    newTopicName: string;
    onNewTopicNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSubmit: (e: React.FormEvent) => Promise<boolean> | boolean;
    onCancel: () => void;
    isSubmitting?: boolean;
    error?: Error | null;
    onClearError?: () => void;
}

export default function CreateTopicForm({
                                            newTopicName,
                                            onNewTopicNameChange,
                                            onSubmit,
                                            onCancel,
                                            isSubmitting = false,
                                            error,
                                            onClearError,
                                        }: CreateTopicFormProps) {

    const handleCancel = () => {
        onCancel();
        onClearError?.();
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onNewTopicNameChange(e);
        if (error && onClearError) {
            onClearError();
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const result = await onSubmit(e);
            if (result === true) {
                onCancel();
            } else {
            }
        } catch (error) {
        }
    };

    return (
        <div className="bg-[var(--background)] border-b border-[var(--border)]">
            {error && (
                <div
                    className="mx-3 mt-3 p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md flex items-center justify-between">
                    <div className="flex items-center gap-2 text-red-800 dark:text-red-200">
                        <AlertCircle className="h-4 w-4"/>
                        <span className="text-sm">{error.message}</span>
                    </div>
                    {onClearError && (
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={onClearError}
                            className="h-6 w-6 p-0 text-red-800 dark:text-red-200 hover:bg-red-100 dark:hover:bg-red-800/20"
                        >
                            <X className="h-3 w-3"/>
                        </Button>
                    )}
                </div>
            )}
            <form onSubmit={handleSubmit} className="p-3">
                <div className="flex gap-2">
                    <Input
                        placeholder="Topic name"
                        value={newTopicName}
                        onChange={handleInputChange}
                        autoFocus
                        disabled={isSubmitting}
                        className="h-9 bg-[var(--background)] border border-[var(--border)]"
                    />
                    <Button
                        type="submit"
                        size="sm"
                        className="h-9 px-2"
                        disabled={isSubmitting || !newTopicName.trim()}
                    >
                        {isSubmitting ? (
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"/>
                        ) : (
                            <Check className="h-4 w-4"/>
                        )}
                    </Button>
                    <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        className="h-9 px-2"
                        onClick={handleCancel}
                        disabled={isSubmitting}
                    >
                        <X className="h-4 w-4"/>
                    </Button>
                </div>
            </form>
        </div>
    );
}