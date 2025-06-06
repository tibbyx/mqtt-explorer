import {useCallback, useState} from "react";
import {apiClient} from "@/api/client.ts";
import {endpoints} from "@/api/endpoints.ts";
import type {Topic} from "@/lib/types.ts";

export function useTopics() {
    const [topics, setTopics] = useState<Topic[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const fetchTopics = useCallback(async () => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await apiClient.get<{ topics: string[] }>(
                endpoints.subscribedTopics
            )

            const mappedTopics: Topic[] = response.topics.map((name: any) => ({
                id: name,
                name: name,
                subscribed: true
            }));

            setTopics(mappedTopics)
            console.log("The Topics are here!", response);
            return response;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, [])

    return {
        topics,
        isLoading,
        error,
        fetchTopics
    }
}