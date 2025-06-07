// api/hooks/useSubscribeTopics.ts
import {useCallback, useState} from "react";
import {apiClient} from "@/api/client.ts";
import {endpoints} from "@/api/endpoints.ts";

interface TopicResult {
    status: string;
    message: string;
}

interface SubscribeTopicsRequest {
    topics: string[];
}

interface SubscribeTopicsResponse {
    result: Record<string, TopicResult>;
}

export function useSubscribeTopics() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const subscribeToTopics = useCallback(async (topics: string[]) => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await apiClient.post<
                SubscribeTopicsResponse,
                SubscribeTopicsRequest
            >(endpoints.subscribeToTopic, {topics});

            console.log("Subscription response:", response);
            return response;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, []);

    return {
        subscribeToTopics,
        isLoading,
        error,
    };
}