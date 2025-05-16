import {formatDistanceToNow} from "date-fns"
import type {Message, QoSLevel} from "@/lib/types.ts";
import {Badge} from "@/components/ui/badge.tsx";

export function MessageItem({message}: { message: Message }) {
    const getQoSBadgeVariant = (q: QoSLevel) =>
        q === 0
            ? "bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-900/20 dark:text-blue-400 dark:border-blue-800/50"
            : q === 1
                ? "bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-900/20 dark:text-amber-400 dark:border-amber-800/50"
                : "bg-[#7a62f6]/10 text-[#7a62f6] border-[#7a62f6]/20 dark:bg-[#7a62f6]/20 dark:border-[#7a62f6]/30";

    return (
        <div
            className="p-4 border-b bg-gray-50 dark:bg-gray-900/50 hover:bg-white dark:hover:bg-gray-950 transition-colors">
            <div className="flex items-start justify-between mb-2">
                <Badge
                    variant="outline"
                    className={`flex items-center gap-1 font-medium ${getQoSBadgeVariant(
                        message.qos
                    )}`}
                >
                    QoS {message.qos}
                </Badge>
                <span className="text-xs text-gray-500 dark:text-gray-400">
          {formatDistanceToNow(new Date(message.timestamp), {
              addSuffix: true,
          })}
        </span>
            </div>
            <div
                className="
          font-mono text-sm bg-white dark:bg-gray-950 p-3 shadow-sm mb-2
          whitespace-pre-wrap break-words
          text-gray-800 dark:text-gray-200
        "
            >
                {message.payload}
            </div>
        </div>
    );
}