// App.tsx
import React, { useEffect, useState } from "react";
import { fetchVideos } from "../services/apiService";

interface Video {
  title: string;
}

const App: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchVideos();
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

  return (
    <div>
      <h1>Videos</h1>
      {error && <p>Error: {error}</p>}
      <ul>
        {videos.map((video, index) => (
          <li key={index}>{video.title}</li>
        ))}
      </ul>
    </div>
  );
};

export default App;
