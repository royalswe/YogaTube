// apiService.ts

// Define a service for API calls
export const fetchData = async (url: string) => {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error fetching videos:", error);
        throw error;
    }
};
