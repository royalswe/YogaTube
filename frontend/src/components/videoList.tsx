// App.tsx
import type { Video } from "../models/video"; // Import the Video interface
import React, { useEffect, useState } from "react";
import { fetchData } from "../services/apiService";

interface AppProps {
  onVideoClick?: (videoId: number) => void;
}

const App: React.FC<AppProps> = ({ onVideoClick }) => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchData("http://localhost:8080/api/v1/videos");
        setVideos(data);
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("An unknown error occurred");
        }
      }
    };

    getVideos();
  }, []);

  const handleVideoClick = (videoId: number) => {
    if (onVideoClick) {
      onVideoClick(videoId);
    } else {
      console.log(`Clicked video with ID: ${videoId}`);
    }
  };

  return (
    <div>
      <h1>Videos</h1>
      {error && <p>Error: {error}</p>}
      <ul>
        {videos.map((video) => (
          <li key={video.id} onClick={() => handleVideoClick(video.id)}>{video.title}</li>
        ))}
      </ul>
    </div>
  );
};

export default App;
