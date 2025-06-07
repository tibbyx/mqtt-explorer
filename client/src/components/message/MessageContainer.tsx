import React from "react";
import type {Message} from "@/lib/types";
import {MessageItem} from "@/components/message/MessageItem";

export interface MessagesContainerProps {
    messages: Message[];
    isLoading?: boolean;
    error?: Error | null;
}

export const MessagesContainer = React.forwardRef<
    HTMLDivElement,
    MessagesContainerProps
>(({messages, isLoading = false, error = null}, ref) => {
    if (isLoading) {
        return (
            <div className="flex-1 flex items-center justify-center">
                <div className="text-center">
                    <div
                        className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100 mx-auto mb-2"></div>
                    <p className="text-sm text-gray-500">Loading messages...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex-1 flex items-center justify-center">
                <div className="text-center text-red-500">
                    <p className="font-medium">Failed to load messages</p>
                    <p className="text-sm">{error.message}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="flex-1 overflow-hidden">
            <div ref={ref} className="h-full overflow-y-auto px-4 py-2 space-y-2">
                {messages.length === 0 ? (
                    <div className="flex items-center justify-center h-full">
                        <p>No messages yet</p>
                    </div>
                ) : (
                    messages.map((message) => (
                        <MessageItem key={message.id} message={message}/>
                    ))
                )}
            </div>
        </div>
    );
});
