import {useState, useCallback} from 'react';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

// body to match backend
interface SubscribeTopicsRequest {
    BrokerUserIds: {
        BrokerId: number;
        UserId: number;
    };
    topics: string[];
}

// result from backend
interface TopicResult {
    status: string;
    message: string;
}

// response
interface SubscribeTopicsResponse {
    result: Record<string, TopicResult>;
}

interface SubscriptionResult {
    topic: string;
    success: boolean;
    status: string;
    message: string;
}

// error response
interface SubscribeError {
    Unauthorized?: string;
    badJson?: string;
    terribleJson?: string;
    InternalServerError?: string;
    Error?: string;
}

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

export function useSubscribeTopics() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const [lastSubscriptionResults, setLastSubscriptionResults] = useState<SubscriptionResult[]>([]);

    const subscribeToTopics = useCallback(async (topics: string[]) => {
        setIsLoading(true);
        setError(null);

        try {
            const brokerId = getFromLocalStorage('brokerId', -1);
            const userId = getFromLocalStorage('userId', -1);

            if (brokerId <= 0 || userId <= 0) {
                throw new Error('Invalid broker ID or user ID. Please check your localStorage settings.');
            }

            if (!Array.isArray(topics) || topics.length === 0) {
                throw new Error('Topics array cannot be empty.');
            }

            const validTopics = topics
                .map(topic => topic.trim())
                .filter(topic => topic.length > 0);

            if (validTopics.length === 0) {
                throw new Error('No valid topics provided.');
            }

            const requestBody: SubscribeTopicsRequest = {
                BrokerUserIds: {
                    BrokerId: brokerId,
                    UserId: userId,
                },
                topics: validTopics,
            };

            const response = await apiClient.post<SubscribeTopicsResponse>(
                endpoints.subscribeToTopic,
                requestBody
            );

            // results into a more usable format
            const results: SubscriptionResult[] = Object.entries(response.result).map(
                ([topic, result]) => ({
                    topic,
                    success: result.status === 'Fine',
                    status: result.status,
                    message: result.message,
                })
            );

            setLastSubscriptionResults(results);

            const allSuccessful = results.every(result => result.success);
            const hasErrors = results.some(result => !result.success);

            return {
                success: allSuccessful,
                hasErrors,
                results,
                successfulTopics: results.filter(r => r.success).map(r => r.topic),
                failedTopics: results.filter(r => !r.success).map(r => r.topic),
            };

        } catch (err) {
            let errorMessage = 'Failed to subscribe to topics';

            if (err instanceof Error) {
                errorMessage = err.message;
            } else if (typeof err === 'object' && err !== null) {
                const backendError = err as any;
                if (backendError.response?.data) {
                    const errorData = backendError.response.data as SubscribeError;
                    errorMessage = errorData.Unauthorized ||
                        errorData.badJson ||
                        errorData.terribleJson ||
                        errorData.InternalServerError ||
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

    const subscribeToTopic = useCallback(async (topic: string) => {
        return subscribeToTopics([topic]);
    }, [subscribeToTopics]);

    const clearError = useCallback(() => {
        setError(null);
    }, []);

    const clearResults = useCallback(() => {
        setLastSubscriptionResults([]);
    }, []);

    const getCurrentIds = useCallback(() => ({
        brokerId: getFromLocalStorage('brokerId', -1),
        userId: getFromLocalStorage('userId', -1),
    }), []);

    // get results for a specific topic
    const getTopicResult = useCallback((topic: string) => {
        return lastSubscriptionResults.find(result => result.topic === topic);
    }, [lastSubscriptionResults]);

    // check if a topic was successfully subscribed
    const isTopicSubscribed = useCallback((topic: string) => {
        const result = getTopicResult(topic);
        return result?.success ?? false;
    }, [getTopicResult]);

    return {
        subscribeToTopics,
        subscribeToTopic,
        isLoading,
        error,
        lastSubscriptionResults,
        clearError,
        clearResults,
        getCurrentIds,
        getTopicResult,
        isTopicSubscribed,
    };
}