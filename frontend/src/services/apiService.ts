// apiService.ts

// Define a service for API calls
export const fetchVideos = async () => {
    try {
        const response = await fetch("http://localhost:8080/api/v1/videos");
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
