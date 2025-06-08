import {formatDistanceToNow} from "date-fns"
import type {Message} from "@/lib/types.ts";
import {Badge} from "@/components/ui/badge.tsx";
import React from "react";

export const MessageItem = React.memo(function MessageItem({message}: { message: Message }) {
    return (
        <div className="p-4 bg-[var(--background)] border border-[var(--border)]">
            <div className="flex items-start justify-between mb-2">
                <Badge
                    variant="outline"
                    className={"flex items-center gap-1 font-medium"}
                >
                    {message.ClientId}
                </Badge>
                <span className="text-xs">
                    {formatDistanceToNow(new Date(message.timestamp), {
                        addSuffix: true,
                    })}
                </span>
            </div>
            <div
                className="p-3 mb-2 whitespace-pre-wrap break-words bg-white dark:bg-[var(--popover)] border border-[var(--border)]">
                {message.payload}
            </div>
        </div>
    );
});