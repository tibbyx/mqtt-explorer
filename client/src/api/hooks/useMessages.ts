// hooks/useMessages.ts
import {useState, useCallback, useEffect, useRef} from 'react';
import type {Message, QoSLevel} from '@/lib/types';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

interface TopicMessagesResponse {
    topic: string;
    messages: string[];
}

// Helper function to compare message arrays
const messagesEqual = (a: Message[], b: Message[]) => {
    if (a.length !== b.length) return false;
    return a.every((msg, index) => msg.payload === b[index]?.payload);
};

// Helper function to create stable message ID
const createMessageId = (topicName: string, payload: string, index: number) => {
    // Create a simple hash of the payload for stable IDs
    const hash = payload.split('').reduce((a, b) => {
        a = ((a << 5) - a) + b.charCodeAt(0);
        return a & a;
    }, 0);
    return `${topicName}-${index}-${Math.abs(hash)}`;
};

export function useMessages() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const [currentTopic, setCurrentTopic] = useState<string | null>(null);
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [refreshInterval, setRefreshInterval] = useState(1000);

    const intervalRef = useRef<NodeJS.Timeout | null>(null);
    const lastFetchTime = useRef<number>(0);

    const fetchMessages = useCallback(async (topicName: string, silent = false) => {
        if (!silent) {
            setIsLoading(true);
        }
        setError(null);

        try {
            const response = await apiClient.get<TopicMessagesResponse>(
                `${endpoints.getMessageFromTopic}?topic=${encodeURIComponent(topicName)}`
            );

            const currentTime = Date.now();

            // Transform string messages into Message objects with stable IDs
            const transformedMessages: Message[] = response.messages.map((payload, index) => ({
                id: createMessageId(topicName, payload, index),
                topic: response.topic,
                payload: payload,
                timestamp: new Date(lastFetchTime.current || currentTime).toISOString(),
                qos: 0 as QoSLevel,
            }));

            // Only update if messages actually changed
            setMessages(prevMessages => {
                if (messagesEqual(prevMessages, transformedMessages)) {
                    return prevMessages; // Return same reference to prevent re-render
                }
                lastFetchTime.current = currentTime;
                return transformedMessages;
            });

            return transformedMessages;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            throw errorObj;
        } finally {
            if (!silent) {
                setIsLoading(false);
            }
        }
    }, []);

    // Start watching a topic with auto-refresh
    const startWatching = useCallback((topicName: string) => {
        setCurrentTopic(topicName);
        fetchMessages(topicName, false); // First fetch with loading
    }, [fetchMessages]);

    // Stop watching
    const stopWatching = useCallback(() => {
        setCurrentTopic(null);
        setMessages([]);
        setError(null);
        lastFetchTime.current = 0;
        if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
        }
    }, []);

    const toggleAutoRefresh = useCallback(() => {
        setAutoRefresh(prev => !prev);
    }, []);

    const setRefreshRate = useCallback((milliseconds: number) => {
        setRefreshInterval(milliseconds);
    }, []);

    const refresh = useCallback(() => {
        if (currentTopic) {
            fetchMessages(currentTopic, false); // Manual refresh with loading
        }
    }, [currentTopic, fetchMessages]);

    // Auto-refresh effect
    useEffect(() => {
        if (currentTopic && autoRefresh) {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }

            intervalRef.current = setInterval(() => {
                fetchMessages(currentTopic, true); // Silent refresh
            }, refreshInterval);

            return () => {
                if (intervalRef.current) {
                    clearInterval(intervalRef.current);
                    intervalRef.current = null;
                }
            };
        } else {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
                intervalRef.current = null;
            }
        }
    }, [currentTopic, autoRefresh, refreshInterval, fetchMessages]);

    useEffect(() => {
        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
        };
    }, []);

    const clearMessages = useCallback(() => {
        setMessages([]);
        setError(null);
        lastFetchTime.current = 0;
    }, []);

    const addMessage = useCallback((message: Message) => {
        setMessages(prev => [...prev, message]);
        lastFetchTime.current = Date.now();
    }, []);

    return {
        messages,
        isLoading,
        error,
        currentTopic,
        autoRefresh,
        refreshInterval,
        fetchMessages,
        clearMessages,
        addMessage,
        startWatching,
        stopWatching,
        refresh,
        toggleAutoRefresh,
        setRefreshRate,
    };
}