import React, {useState, useRef, useEffect} from "react";
import type {Message, QoSLevel, Topic} from "@/lib/types";
import {NoTopicSelectedView} from "./NoTopicSelectedView";
import {MessageTopicHeader} from "./MessageTopicHeader.tsx";
import {MessagesContainer} from "./MessageContainer";
import {MessageComposer} from "./MessageComposer";

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
    const [publishQosLevel, setPublishQosLevel] = useState<QoSLevel>(0);

    const [filterQos, setFilterQos] = useState<QoSLevel | null>(null);
    const [isSubscribed, setIsSubscribed] = useState<boolean | null>(null);
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setIsSubscribed(topic?.subscribed ?? null);
    }, [topic]);

    const filteredMessages = messages.filter(
        (message) => filterQos === null || message.qos === filterQos
    );

    useEffect(() => {
        if (filterQos === null && messagesContainerRef.current) {
            const container = messagesContainerRef.current;
            container.scrollTop = container.scrollHeight;
        }
    }, [filteredMessages.length, filterQos]);

    const handlePublish = (e: React.FormEvent) => {
        e.preventDefault();
        if (!topic || !messageText.trim()) return;

        onPublish(topic.name, messageText.trim(), publishQosLevel);
        setMessageText("");
    };

    const handleSubscriptionToggle = () => {
        if (!topic) return;

        const newSubscriptionState = !isSubscribed;
        if (newSubscriptionState) {
            onSubscribe?.(topic.id);
        } else {
            onUnsubscribe?.(topic.id);
        }
        setIsSubscribed(newSubscriptionState);
    };

    const handleQosFilterChange = (value: string) =>
        setFilterQos(value === "any" ? null : (parseInt(value) as QoSLevel));

    if (!topic) {
        return <NoTopicSelectedView/>;
    }

    return (
        <div className="flex-1 flex flex-col h-full bg-gray-50 dark:bg-gray-900">
            <MessageTopicHeader
                topicName={topic.name}
                isSubscribed={!!isSubscribed}
                onSubscriptionToggle={handleSubscriptionToggle}
                filterQos={filterQos}
                onFilterChange={handleQosFilterChange}
            />

            <MessagesContainer
                ref={messagesContainerRef}
                messages={filteredMessages}
            />

            <MessageComposer
                messageText={messageText}
                onMessageChange={setMessageText}
                qosLevel={publishQosLevel}
                onQosChange={setPublishQosLevel}
                onPublish={handlePublish}
            />
        </div>
    );
}
