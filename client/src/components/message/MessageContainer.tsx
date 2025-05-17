import React from "react";
import type {Message} from "@/lib/types";
import {MessageItem} from "@/components/message/MessageItem";

export interface MessagesContainerProps {
    messages: Message[];
}

export const MessagesContainer = React.forwardRef<
    HTMLDivElement,
    MessagesContainerProps
>(({messages}, ref) => {
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

MessagesContainer.displayName = "MessagesContainer";
