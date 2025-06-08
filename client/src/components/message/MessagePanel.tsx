import React, {useState, useRef, useEffect} from "react";
import type {Message, QoSLevel, Topic} from "@/lib/types";
import {NoTopicSelectedView} from "./NoTopicSelectedView";
import {MessageTopicHeader} from "./MessageTopicHeader.tsx";
import {MessagesContainer} from "./MessageContainer";
import {MessageComposer} from "./MessageComposer";
import {useMessages} from "@/api/hooks/useMessages.ts";
import {useSendMessage} from "@/api/hooks/useSendMessage.ts";

export interface MessagePanelProps {
    topic: Topic | null;
    messages: Message[];
    onSubscribe?: (topicId: string) => void;
    onUnsubscribe?: (topicId: string) => void;
    isLoading?: boolean;
    error?: Error | null;
}

export function MessagePanel({
                                 topic,
                                 onSubscribe,
                                 onUnsubscribe,
                             }: MessagePanelProps) {
    const {
        messages,
        isLoading,
        error,
        startWatching,
        stopWatching,
        refresh,
    } = useMessages();

    const {
        sendMessage,
        isLoading: isSending,
        error: sendError,
        clearError: clearSendError,
    } = useSendMessage();

    const [messageText, setMessageText] = useState("");
    const [publishQosLevel, setPublishQosLevel] = useState<QoSLevel>(0);
    const [filterQos, setFilterQos] = useState<QoSLevel | null>(null);
    const [isSubscribed, setIsSubscribed] = useState<boolean | null>(null);
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (topic) {
            startWatching(topic.Topic);
        } else {
            stopWatching();
        }
    }, [topic?.Topic, startWatching, stopWatching]);

    const filteredMessages = messages.filter(
        (message) => filterQos === null || message.qos === filterQos
    );

    useEffect(() => {
        if (filterQos === null && messagesContainerRef.current) {
            const container = messagesContainerRef.current;
            container.scrollTop = container.scrollHeight;
        }
    }, [filteredMessages.length, filterQos]);

    // handlePublish to use the sendMessage hook
    const handlePublish = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!topic || !messageText.trim()) return;
        try {
            clearSendError();
            await sendMessage(topic.Topic, messageText.trim());
            setMessageText("");
            setTimeout(() => refresh(), 500);
        } catch (error) {
            console.error('Failed to send message:', error);
        }
    };

    const handleSubscriptionToggle = () => {
        if (!topic) return;
        const newSubscriptionState = !isSubscribed;
        if (newSubscriptionState) {
            onSubscribe?.(topic.Id);
            startWatching(topic.Topic);
        } else {
            onUnsubscribe?.(topic.Id);
            stopWatching();
        }
        setIsSubscribed(newSubscriptionState);
    };

    const handleQosFilterChange = (value: string) =>
        setFilterQos(value === "any" ? null : (parseInt(value) as QoSLevel));

    if (!topic) {
        return <NoTopicSelectedView/>;
    }

    return (
        <div className="flex-1 flex flex-col h-full bg-gray-50 dark:bg-[var(--secondary-foreground)]">
            <MessageTopicHeader
                topicName={topic.Topic}
                isSubscribed={!!isSubscribed}
                onSubscriptionToggle={handleSubscriptionToggle}
                filterQos={filterQos}
                onFilterChange={handleQosFilterChange}
                isLoading={isLoading}
            />
            <MessagesContainer
                ref={messagesContainerRef}
                messages={filteredMessages}
                isLoading={isLoading}
                error={error || sendError}
            />
            <MessageComposer
                messageText={messageText}
                onMessageChange={setMessageText}
                qosLevel={publishQosLevel}
                onQosChange={setPublishQosLevel}
                onPublish={handlePublish}
                isPublishing={isSending}
                publishError={sendError}
                onClearError={clearSendError}
            />
        </div>
    );
}