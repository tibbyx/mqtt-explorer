import React, {useState, useRef, useEffect} from "react"
import {Send, MessageSquare} from "lucide-react"
import type {Message, QoSLevel, Topic} from "@/lib/types"
import {Button} from "../ui/button"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "../ui/select"
import {MessageItem} from "@/components/message/MessageItem.tsx";
import {Textarea} from "../ui/textarea"
import {Label} from "../ui/label"

export interface MessagePanelProps {
    topic: Topic | null;
    messages: Message[];
    onPublish: (topic: string, payload: string, qos: QoSLevel) => void;
    onSubscribe?: (topicId: string) => void;
    onUnsubscribe?: (topicId: string) => void;
}

export function MessagePanel({
                                 topic,
                                 messages,
                                 onPublish,
                                 onSubscribe,
                                 onUnsubscribe,
                             }: MessagePanelProps) {
    const [messageText, setMessageText] = useState("");
    const [qosLevel, setQosLevel] = useState<QoSLevel>(0);
    const [filterQos, setFilterQos] = useState<QoSLevel | null>(null);
    const [localTopicSubscribed, setLocalTopicSubscribed] =
        useState<boolean | null>(null);

    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setLocalTopicSubscribed(topic ? topic.subscribed : null);
    }, [topic]);

    const filteredMessages = messages.filter(
        (m) => filterQos === null || m.qos === filterQos
    );
    const isFiltering = filterQos !== null;

    useEffect(() => {
        if (!isFiltering && messagesContainerRef.current) {
            const el = messagesContainerRef.current;
            el.scrollTop = el.scrollHeight;
        }
    }, [filteredMessages.length, isFiltering]);

    const handlePublish = (e: React.FormEvent) => {
        e.preventDefault();
        if (topic && messageText.trim()) {
            onPublish(topic.name, messageText.trim(), qosLevel);
            setMessageText("");
        }
    };

    const handleSubscriptionToggle = () => {
        if (!topic) return;
        if (localTopicSubscribed) {
            onUnsubscribe?.(topic.id);
            setLocalTopicSubscribed(false);
        } else {
            onSubscribe?.(topic.id);
            setLocalTopicSubscribed(true);
        }
    };

    const handleQosFilterChange = (value: string) =>
        setFilterQos(value === "any" ? null : (parseInt(value) as QoSLevel));

    if (!topic) {
        return (
            <div className="flex-1 flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <div className="text-center p-8 max-w-md">
                    <MessageSquare className="h-16 w-16 mx-auto mb-4 text-gray-300 dark:text-gray-700"/>
                    <h3 className="text-xl font-medium mb-2 text-gray-700 dark:text-gray-300">
                        No topic selected
                    </h3>
                    <p className="text-gray-500 dark:text-gray-400">
                        Select a topic from the list to view messages or create a new topic
                        to get started.
                    </p>
                </div>
            </div>
        );
    }

    const isSubscribed =
        localTopicSubscribed !== null ? localTopicSubscribed : topic.subscribed;

    return (
        <div className="flex-1 flex flex-col h-full bg-gray-50 dark:bg-gray-900">
            {/* Header */}
            <div
                className="p-4 border-b border-gray-200 dark:border-gray-800 flex items-center justify-between bg-white dark:bg-gray-950">
                <div className="flex items-center gap-2">
                    <h2 className="font-medium text-gray-700 dark:text-gray-200">
                        {topic.name}
                    </h2>
                    <Button
                        size="sm"
                        variant={isSubscribed ? "outline" : "default"}
                        onClick={handleSubscriptionToggle}
                        className={
                            isSubscribed
                                ? "border-[#7a62f6] text-[#7a62f6] hover:bg-[#7a62f6]/10"
                                : "bg-[#7a62f6] hover:bg-[#6952e3] text-white"
                        }
                    >
                        {isSubscribed ? "Unsubscribe" : "Subscribe"}
                    </Button>
                </div>
                <Select
                    value={filterQos === null ? "any" : filterQos.toString()}
                    onValueChange={handleQosFilterChange}
                >
                    <SelectTrigger
                        className="w-28 h-9 rounded-lg border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 text-gray-700 dark:text-gray-200 focus:ring-[#7a62f6]">
                        <SelectValue placeholder="QoS" className="font-medium"/>
                    </SelectTrigger>
                    <SelectContent className="rounded-lg border-gray-200 dark:border-gray-800">
                        <SelectItem value="any">Any QoS</SelectItem>
                        <SelectItem value="0">QoS 0</SelectItem>
                        <SelectItem value="1">QoS 1</SelectItem>
                        <SelectItem value="2">QoS 2</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            <div className="flex-1 overflow-hidden">
                <div
                    ref={messagesContainerRef}
                    className="h-full overflow-y-auto px-4 py-2 space-y-2"
                >
                    {filteredMessages.length === 0 ? (
                        <div className="flex items-center justify-center h-full">
                            <p className="text-gray-500 dark:text-gray-400">No messages yet</p>
                        </div>
                    ) : (
                        filteredMessages.map((m) => (
                            <MessageItem key={m.id} message={m}/>
                        ))
                    )}
                </div>
            </div>

            <div className="p-4 border-t border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950">
                <form onSubmit={handlePublish} className="space-y-3">
                    <div className="flex gap-4">
                        <div className="flex-1">
                            <Textarea
                                placeholder="Enter message payload..."
                                value={messageText}
                                onChange={(e) => setMessageText(e.target.value)}
                                className="min-h-[80px] rounded-lg border-gray-200 dark:border-gray-800 focus-visible:ring-[#7a62f6] text-gray-800 dark:text-gray-200"
                            />
                        </div>
                        <div className="space-y-4">
                            <div className="space-y-2">
                                <Label
                                    htmlFor="qos-select"
                                    className="text-gray-700 dark:text-gray-300"
                                >
                                    QoS Level
                                </Label>
                                <Select
                                    value={qosLevel.toString()}
                                    onValueChange={(v) =>
                                        setQosLevel(parseInt(v) as QoSLevel)
                                    }
                                >
                                    <SelectTrigger
                                        id="qos-select"
                                        className="w-28 rounded-lg border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 text-gray-700 dark:text-gray-200 focus:ring-[#7a62f6]"
                                    >
                                        <SelectValue placeholder="QoS" className="font-medium"/>
                                    </SelectTrigger>
                                    <SelectContent className="rounded-lg border-gray-200 dark:border-gray-800">
                                        <SelectItem value="0">QoS 0</SelectItem>
                                        <SelectItem value="1">QoS 1</SelectItem>
                                        <SelectItem value="2">QoS 2</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                    <div className="flex justify-end">
                        <Button
                            type="submit"
                            disabled={!messageText.trim()}
                            className="bg-[#7a62f6] hover:bg-[#6952e3] text-white"
                        >
                            <Send className="h-4 w-4 mr-2"/>
                            Publish
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}