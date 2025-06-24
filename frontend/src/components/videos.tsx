// App.tsx
import React, { useEffect, useState } from "react";
import { fetchData } from "../services/apiService";

interface Thumbnail {
  url: string;
  width: number;
  height: number;
}

interface Thumbnails {
  default?: Thumbnail;
  medium?: Thumbnail;
  [key: string]: Thumbnail | undefined;
}

interface ResourceId {
  kind?: string;
  videoId?: string;
}

interface Video {
  title: string;
  publishedAt?: string;
  description?: string;
  videoOwnerChannelTitle?: string;
  thumbnails?: Thumbnails;
  resourceId?: ResourceId;
}

const App: React.FC = () => {
  const [video, setVideos] = useState<Video>();
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState<number>(0);
  const apiDomain = import.meta.env.VITE_API_DOMAIN ?? "http://localhost:83080";
  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchData(`${apiDomain}/api/v1/video?offset=${offset}`);
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
  }, [offset, apiDomain]);

  const handleNext = () => {
    setOffset((prevOffset) => prevOffset + 1);
  };

  const handlePrevious = () => {
    setOffset((prevOffset) => Math.max(prevOffset - 1, 0));
  };

  return (
    <div>
      {error && <p>Error: {error}</p>}
      {video && <h2>{video.title}</h2>}
      {video && (
        <div>
          <p><strong>Description:</strong> {video.description?.slice(0, 100)}...</p>
          <p><strong>Channel Title:</strong> {video.videoOwnerChannelTitle}</p>
          <div>
            <strong>Thumbnails:</strong>
            <div>
              <p>Default: {video.thumbnails?.default?.width}</p>
              <img src={video.thumbnails?.default?.url} alt="Default Thumbnail" />
            </div>
          </div>
          <div>
            <strong>Resource ID:</strong>
            <p>Kind: {video.resourceId?.kind}</p>
            <p>Video ID: {video.resourceId?.videoId}</p>
          </div>
        </div>
      )}

    {video?.resourceId?.videoId && (
      <div>
        <h3>Watch Video:</h3>
        <iframe
      width="560"
      height="315"
      src={`https://www.youtube.com/embed/${video.resourceId.videoId}`}
      title="YouTube video player"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowFullScreen
        ></iframe>
      </div>
    )}

    <div>
      <button onClick={handlePrevious} disabled={offset === 0}>Previous</button>
      <button onClick={handleNext}>Next</button>
    </div>

    </div>
    
  );
};

export default App;
