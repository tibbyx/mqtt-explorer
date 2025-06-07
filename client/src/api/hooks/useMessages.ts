// hooks/useMessages.ts
import {useState, useCallback} from 'react';
import type {Message, QoSLevel} from '@/lib/types';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

// Backend response type based on your Go handler
interface TopicMessagesResponse {
    topic: string;
    messages: string[];
}

export function useMessages() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const fetchMessages = useCallback(async (topicName: string) => {
        setIsLoading(true);
        setError(null);

        try {
            // Call your backend endpoint with topic query parameter
            const response = await apiClient.get<TopicMessagesResponse>(
                `${endpoints.getMessageFromTopic}?topic=${encodeURIComponent(topicName)}`
            );

            // Transform string messages into Message objects
            const transformedMessages: Message[] = response.messages.map((payload, index) => ({
                id: `${topicName}-${Date.now()}-${index}`, // Generate unique ID
                topic: response.topic,
                payload: payload,
                timestamp: new Date().toISOString(), // Use current time since backend doesn't provide it
                qos: 0 as QoSLevel, // Default QoS since backend doesn't provide it
            }));

            setMessages(transformedMessages);
            return transformedMessages;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const clearMessages = useCallback(() => {
        setMessages([]);
        setError(null);
    }, []);

    const addMessage = useCallback((message: Message) => {
        setMessages(prev => [...prev, message]);
    }, []);

    return {
        messages,
        isLoading,
        error,
        fetchMessages,
        clearMessages,
        addMessage, // For real-time message updates
    };
}