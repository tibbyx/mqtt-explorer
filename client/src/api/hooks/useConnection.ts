import type {Credentials} from "@/lib/types.ts";
import {apiClient} from "@/api/client.ts";
import {endpoints} from "@/api/endpoints.ts";
import {useState} from "react";

export function useConnection() {
    const [error, setError] = useState<Error | null>(null);

    const connect = async (credentials: Credentials) => {
        setError(null);

        try {
            const response = await apiClient.post(
                endpoints.credentials,
                credentials
            );
            console.log("The connection was super duper!" + response);
            return response;
        } catch (err) {
            const errorObj = err instanceof Error ? err : new Error(String(err));
            setError(errorObj);
            throw errorObj;
        }
    }

    return {
        connect,
        error,
    };
}

