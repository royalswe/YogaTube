// App.tsx
import type { Video } from "../models/video"; // Import the Video interface
import React, { useEffect, useState, useRef } from "react";
import { fetchData } from "../services/apiService";
import "./videos.css"; // Add a CSS file for styling
import VideoList from "./videoList";

const App: React.FC = () => {
  const [video, setVideos] = useState<Video>();
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState<number>(0);
  const [showFullDescription, setShowFullDescription] = useState<boolean>(false);
  const [showVideoList, setShowVideoList] = useState<boolean>(false);
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const videoListRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchData(`/api/v1/video?offset=${offset}`);
        if (data?.exceeded) {
          setOffset(0)
          alert("No more videos available. Resetting to the first video.");
          return;
        }
        setVideos(data);
      } catch (err) {
        setOffset(0)
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

  const handleShowVideoList = () => {
    setShowVideoList(true);
    videoListRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const handleVideoClick = (clickedId: number) => {
    if (video?.id !== undefined) {
      const newOffset = offset + (clickedId - video.id);
      setOffset(newOffset);
      setShowVideoList(false);
    }
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
          <div className="video-info">
            <p><strong>Channel Title:</strong> {video.videoOwnerChannelTitle}</p>
            <div
              className="description"
              onClick={() => setShowFullDescription((prev) => !prev)}
            >
              <strong>Description:</strong>{" "}
              {showFullDescription || !video.description || video.description.length <= 70 ? (
                video.description
              ) : (
                <>
                  {video.description.slice(0, 70)}...
                  <span
                    style={{ fontWeight: "700", cursor: "pointer", marginLeft: 4 }}
                    onClick={(e) => {
                      e.stopPropagation();
                      setShowFullDescription(true);
                    }}
                  >
                    Show full description ⇩
                  </span>
                </>
              )}
              {showFullDescription && video.description && video.description.length > 70 && (
                <span
                  style={{ fontWeight: "700", cursor: "pointer", marginLeft: 4 }}
                  onClick={(e) => {
                    e.stopPropagation();
                    setShowFullDescription(false);
                  }}
                >
                  Show less
                </span>
              )}
            </div>
          </div>
          <p className="hint">Press 'F' to toggle fullscreen mode for the video.</p>
        </div>
      )}

      <div className="navigation-buttons">
        <button onClick={handlePrevious} disabled={offset === 0} className="button">Previous</button>
        <button onClick={handleNext} className="button">Next</button>
      </div>

      <button onClick={handleShowVideoList} className="button scroll-button">View All Videos</button>

      {showVideoList && (
        <div ref={videoListRef} className="video-list-section">
          <VideoList onVideoClick={handleVideoClick} />
        </div>
      )}
    </div>
  );
};

export default App;
