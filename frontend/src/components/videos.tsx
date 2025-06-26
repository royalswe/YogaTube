// App.tsx
import React, { useEffect, useState, useRef } from "react";
import { fetchData } from "../services/apiService";
import "./videos.css"; // Add a CSS file for styling

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
  const iframeRef = useRef<HTMLIFrameElement>(null);

  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchData(`/api/v1/video?offset=${offset}`);
        if (data?.exceeded) {
          setOffset(0) ; // Reset offset if exceeded
        }
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
  }, [offset]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "f" && iframeRef.current) {
        if (document.fullscreenElement) {
          document.exitFullscreen();
        } else {
          iframeRef.current.requestFullscreen();
        }
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  const handleNext = () => {
    setOffset((prevOffset) => prevOffset + 1);
  };

  const handlePrevious = () => {
    setOffset((prevOffset) => Math.max(prevOffset - 1, 0));
  };

  return (
    <div className="video-page">
      {error && <p className="error">Error: {error}</p>}
      {video && (
        <div className="video-container">
          <h2 className="video-title">{video.title}</h2>
          <iframe
            ref={iframeRef}
            className="video-player"
            width="560"
            height="315"
            src={`https://www.youtube.com/embed/${video.resourceId?.videoId}`}
            title="YouTube video player"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
          ></iframe>
          <p className="hint">Press 'F' to toggle fullscreen mode for the video.</p>
          <div className="video-info">
            <p><strong>Description:</strong> {video.description?.slice(0, 100)}...</p>
            <p><strong>Channel Title:</strong> {video.videoOwnerChannelTitle}</p>
            <div className="thumbnails">
              <strong>Thumbnails:</strong>
              <img src={video.thumbnails?.default?.url} alt="Default Thumbnail" className="thumbnail" />
            </div>
          </div>
        </div>
      )}

      <div className="navigation-buttons">
        <button onClick={handlePrevious} disabled={offset === 0} className="button">Previous</button>
        <button onClick={handleNext} className="button">Next</button>
      </div>
    </div>
  );
};

export default App;
