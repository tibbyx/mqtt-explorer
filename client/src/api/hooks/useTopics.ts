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
            const userId = localStorage.getItem("userId");
            const brokerId = localStorage.getItem("brokerId");

            if (!userId || !brokerId) {
                throw new Error("Missing userId or brokerId in localStorage.");
            }

            const payload = {
                UserId: Number(userId),
                BrokerId: Number(brokerId),
            };

            const response = await apiClient.post<{ Topics: Topic[] }>(
                endpoints.getAllTopics,
                payload
            );

            const topicList = response.Topics;
            setTopics(topicList);
            console.log("Successfully fetched topics:", response);
            return topicList;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            console.error("Error fetching topics:", errorObj);
            throw errorObj;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const addTopic = useCallback((topic: Topic) => {
        setTopics(prev => [...prev, topic]);
    }, []);

    return {
        topics,
        isLoading,
        error,
        fetchTopics,
        addTopic,
    };
}
