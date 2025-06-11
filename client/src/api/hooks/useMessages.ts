import {useState, useCallback, useEffect, useRef} from 'react';
import type {Message, QoSLevel} from '@/lib/types';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

interface SelectMessage {
    Id: number;
    UserId: number;
    ClientId: string;
    TopicId: number;
    BrokerId: number;
    QoS: number;
    Message: string;
    CreationDate: string;
}

interface TopicMessagesResponse {
    topic: string;
    messages: SelectMessage[]
}

interface TopicMessagesRequest {
    brokerUserIDs: {
        brokerId: number;
        userId: number;
    };
    topic: string;
    index: number;
}

// compare message arrays
const messagesEqual = (a: Message[], b: Message[]) => {
    if (a.length !== b.length) return false;
    return a.every((msg, index) =>
        msg.payload === b[index]?.payload &&
        msg.id === b[index]?.id
    );
};

// get values from localStorage
const getFromLocalStorage = (key: string, defaultValue: number): number => {
    try {
        const value = localStorage.getItem(key);
        if (value === null) return defaultValue;
        const parsed = parseInt(value, 10);
        return isNaN(parsed) ? defaultValue : parsed;
    } catch {
        return defaultValue;
    }
};

export function useMessages() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const [currentTopic, setCurrentTopic] = useState<string | null>(null);
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [refreshInterval, setRefreshInterval] = useState(1000);
    const intervalRef = useRef<NodeJS.Timeout | null>(null);
    const fetchMessages = useCallback(async (
        topicName: string,
        silent = false,
        messageIndex = -1
    ) => {
        if (!silent) {
            setIsLoading(true);
        }
        setError(null);

        try {
            const brokerId = getFromLocalStorage('brokerId', -1);
            const userId = getFromLocalStorage('userId', -1);

            if (brokerId < 0 || userId < 0) {
                throw new Error('Invalid broker ID or user ID. Please check your localStorage settings.');
            }

            const requestBody: TopicMessagesRequest = {
                brokerUserIDs: {
                    brokerId,
                    userId,
                },
                topic: topicName,
                index: messageIndex,
            };

            const response = await apiClient.post<TopicMessagesResponse>(
                endpoints.getMessageFromTopic,
                requestBody
            );

            const transformedMessages: Message[] = response.messages.map((dbMessage) => ({
                id: `${dbMessage.Id}-${dbMessage.BrokerId}-${dbMessage.TopicId}`,
                topic: response.topic,
                payload: dbMessage.Message,
                timestamp: dbMessage.CreationDate,
                qos: dbMessage.QoS as QoSLevel,
                ClientId: dbMessage.ClientId

            })).reverse();

            // update if messages actually changed
            setMessages(prevMessages => {
                if (messagesEqual(prevMessages, transformedMessages)) {
                    return prevMessages;
                }
                console.log("New Message: ", transformedMessages);
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

    // watching a topic with auto-refresh
    const startWatching = useCallback((topicName: string) => {
        setCurrentTopic(topicName);
        fetchMessages(topicName, false);
    }, [fetchMessages]);

    // stop watching
    const stopWatching = useCallback(() => {
        setCurrentTopic(null);
        setMessages([]);
        setError(null);
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
            fetchMessages(currentTopic, false);
        }
    }, [currentTopic, fetchMessages]);

    // auto-refresh effect
    useEffect(() => {
        if (currentTopic && autoRefresh) {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
            intervalRef.current = setInterval(() => {
                fetchMessages(currentTopic, true);
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
    }, []);

    const addMessage = useCallback((message: Message) => {
        setMessages(prev => [...prev, message]);
    }, []);

    // current broker and user IDs from localStorage (for display purposes)
    const getCurrentIds = useCallback(() => ({
        brokerId: getFromLocalStorage('brokerId', -1),
        userId: getFromLocalStorage('userId', -1),
    }), []);

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
        getCurrentIds,
    };
}