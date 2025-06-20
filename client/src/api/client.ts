const API_BASE_URL = "http://localhost:3000";

export async function fetchApi<T>(
    endpoint: string,
    options?: RequestInit
): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...(options?.headers || {}),
            },
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API Error ${response.status}: ${errorText || response.statusText}`);
        }

        const data = await response.json();
        return data as T;
    } catch (error) {
        console.error(`Error fetching ${endpoint}:`, error);
        throw error;
    }
}

export const apiClient = {
    get: <T>(endpoint: string, params?: Record<string, string>): Promise<T> => {
        const url = params
            ? `${endpoint}?${new URLSearchParams(params)}`
            : endpoint;

        return fetchApi<T>(url, {
            headers: {
                'Content-Type': 'application/json',
            }
        });
    },

    post: <T, D = any>(endpoint: string, data?: D): Promise<T> => {
        return fetchApi<T>(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: data ? JSON.stringify(data) : undefined,
        });
    },

    /*
    If someone want to add some HTTP -> PUT and DELETE

        put: <T, D = any>(endpoint: string, data?: D): Promise<T> => {
            return fetchApi<T>(endpoint, {
                method: 'PUT',
                body: data ? JSON.stringify(data) : undefined,
            });
        },

        delete: <T>(endpoint: string): Promise<T> => {
            return fetchApi<T>(endpoint, {
                method: 'DELETE',
            });
        },

     */
};