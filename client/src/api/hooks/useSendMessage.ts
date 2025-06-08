import {useState, useCallback} from 'react';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

interface SendMessageRequest {
    BrokerUserIds: {
        BrokerId: number;
        UserId: number;
    };
    Topic: string;
    Message: string;
}

// response
interface SendMessageResponse {
    goodJson: string;
}

// error response
interface SendMessageError {
    Unauthorized?: string;
    badJson?: string;
    "Internal Server Error"?: string;
    Error?: string;
}

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

export function useSendMessage() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const [lastSentMessage, setLastSentMessage] = useState<{
        topic: string;
        message: string;
        timestamp: Date;
    } | null>(null);

    const sendMessage = useCallback(async (topic: string, message: string) => {
        setIsLoading(true);
        setError(null);

        try {
            const brokerId = getFromLocalStorage('brokerId', -1);
            const userId = getFromLocalStorage('userId', -1);

            if (brokerId < 0 || userId < 0) {
                throw new Error('Invalid broker ID or user ID. Please check your localStorage settings.');
            }

            if (!topic.trim()) {
                throw new Error('Topic cannot be empty.');
            }

            if (!message.trim()) {
                throw new Error('Message cannot be empty.');
            }

            const requestBody: SendMessageRequest = {
                BrokerUserIds: {
                    BrokerId: brokerId,
                    UserId: userId,
                },
                Topic: topic.trim(),
                Message: message.trim(),
            };

            const response = await apiClient.post<SendMessageResponse>(
                endpoints.sendMessageToTopic,
                requestBody
            );

            // information about the last sent message
            setLastSentMessage({
                topic: topic.trim(),
                message: message.trim(),
                timestamp: new Date(),
            });

            return {
                success: true,
                message: response.goodJson,
                data: {
                    topic: topic.trim(),
                    message: message.trim(),
                    timestamp: new Date(),
                }
            };

        } catch (err) {
            let errorMessage = 'Failed to send message';

            if (err instanceof Error) {
                errorMessage = err.message;
            } else if (typeof err === 'object' && err !== null) {
                // Handle backend error responses
                const backendError = err as any;
                if (backendError.response?.data) {
                    const errorData = backendError.response.data as SendMessageError;
                    errorMessage = errorData.Unauthorized ||
                        errorData.badJson ||
                        errorData["Internal Server Error"] ||
                        errorData.Error ||
                        'Unknown server error';
                }
            }

            const errorObj = new Error(errorMessage);
            setError(errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const clearError = useCallback(() => {
        setError(null);
    }, []);

    const clearLastSentMessage = useCallback(() => {
        setLastSentMessage(null);
    }, []);

    const getCurrentIds = useCallback(() => ({
        brokerId: getFromLocalStorage('brokerId', -1),
        userId: getFromLocalStorage('userId', -1),
    }), []);

    return {
        sendMessage,
        isLoading,
        error,
        lastSentMessage,
        clearError,
        clearLastSentMessage,
        getCurrentIds,
    };
}