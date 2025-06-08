import {useState, useCallback} from 'react';
import {apiClient} from '@/api/client';
import {endpoints} from '@/api/endpoints';

interface BrokerUserIDs {
    BrokerId: number;
    UserId: number;
}

interface TopicsWrapper {
    BrokerUserIDs: BrokerUserIDs;
    Topics: string[];
}

interface TopicResult {
    Status: string;
    Message: string;
}

interface TopicSubscriptionResponse {
    result: Record<string, TopicResult>;
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

export function useTopicSubscription() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const toggleSubscription = useCallback(async (
        topic: string,
        currentlySubscribed: boolean
    ): Promise<boolean> => {
        setIsLoading(true);
        setError(null);

        try {
            const brokerId = getFromLocalStorage('brokerId', -1);
            const userId = getFromLocalStorage('userId', -1);

            console.log('Toggle subscription for:', {topic, currentlySubscribed, brokerId, userId});

            if (brokerId < 0 || userId < 0) {
                throw new Error('Invalid broker ID or user ID. Please check your localStorage settings.');
            }

            const endpoint = currentlySubscribed
                ? endpoints.unsubscribeToTopic
                : endpoints.subscribeToTopic;

            const requestBody: TopicsWrapper = {
                BrokerUserIDs: {
                    BrokerId: brokerId,
                    UserId: userId,
                },
                Topics: [topic],
            };

            console.log('API Request:', {endpoint, requestBody});

            const response = await apiClient.post<TopicSubscriptionResponse>(
                endpoint,
                requestBody
            );

            console.log('API Response:', response);

            if (!response.result) {
                console.warn('No result object in response');
                return !currentlySubscribed;
            }

            if (response.result[topic]) {
                const topicResult = response.result[topic];
                console.log('Topic result:', topicResult);

                if (topicResult.Status !== 'Fine') {
                    const errorMessage = topicResult.Message || 'No error message provided';
                    throw new Error(`Operation failed: ${errorMessage} (Status: ${topicResult.Status})`);
                }
            } else {
                console.warn('No result found for topic:', topic);
                console.warn('Available topics in result:', Object.keys(response.result));
            }

            const newState = !currentlySubscribed;
            console.log('Subscription toggled successfully, new state:', newState);
            return newState;
        } catch (err) {
            console.error('Error in toggleSubscription:', err);

            let errorMessage = 'Unknown error occurred';

            if (err instanceof Error) {
                errorMessage = err.message || 'Error without message';
            } else if (typeof err === 'string') {
                errorMessage = err;
            } else if (err && typeof err === 'object') {
                if ('response' in err && err.response) {
                    const response = err.response as any;
                    if (response.data) {
                        errorMessage = response.data.Unauthorized ||
                            response.data.terribleJson ||
                            response.data.InternalServerError ||
                            JSON.stringify(response.data);
                    } else {
                        errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                    }
                } else {
                    errorMessage = JSON.stringify(err);
                }
            }

            const errorObj = new Error(errorMessage);
            setError(errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const subscribeToTopic = useCallback(async (topic: string): Promise<void> => {
        await toggleSubscription(topic, false);
    }, [toggleSubscription]);

    const unsubscribeFromTopic = useCallback(async (topic: string): Promise<void> => {
        await toggleSubscription(topic, true);
    }, [toggleSubscription]);

    const getCurrentIds = useCallback(() => ({
        brokerId: getFromLocalStorage('brokerId', -1),
        userId: getFromLocalStorage('userId', -1),
    }), []);

    const clearError = useCallback(() => {
        setError(null);
    }, []);

    return {
        isLoading,
        error,
        toggleSubscription,
        subscribeToTopic,
        unsubscribeFromTopic,
        getCurrentIds,
        clearError,
    };
}