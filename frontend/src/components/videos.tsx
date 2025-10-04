// App.tsx
import type { Video } from "../models/video"; // Import the Video interface
import React, { useEffect, useState, useRef } from "react";
import { YouTubeApi } from "../services/YouTubeApi";
import { fetchData } from "../services/apiService";
import "./videos.css"; // Add a CSS file for styling
import VideoList from "./videoList";

const App: React.FC = () => {
  const [video, setVideos] = useState<Video>();
  const [error, setError] = useState<string | null>(null);
  const [info, setInfo] = useState<string | null>(null);
  const [offset, setOffset] = useState<number>(0);
  const [showFullDescription, setShowFullDescription] = useState<boolean>(false);
  const [showVideoList, setShowVideoList] = useState<boolean>(false);
  const [ytApi, setYtApi] = useState<YouTubeApi | null>(null);
  const videoListRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const getVideos = async () => {
      try {
        const data = await fetchData(`/api/v1/video?offset=${offset}`);
        if (data?.exceeded) {
          // calculate the offset to reset to the first video
          //setOffset(offset - (video?.id || 0));
          setOffset(prevOffset => prevOffset - (video?.id || 0));
          setInfo("No more videos available. Resetting to the first video.");
          setTimeout(() => {
            setInfo(null);
          }, 3000);
          return;
        }
        setVideos(data);
        setError(null);
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
    let api = ytApi;
    const videoId = video?.resourceId?.videoId;
    if (!ytApi && videoId) {
      api = new YouTubeApi("video-player", videoId, () => {
        setYtApi(api);
      });
    } else if (ytApi && videoId) {
      ytApi.loadVideoById(videoId);
    }
    const handleKeyDown = (event: KeyboardEvent) => {
      if (api && !event.ctrlKey && !event.altKey && !event.metaKey && !event.shiftKey) {
        if (event.key === "f") {
          api.toggleFullscreen();
        }
        if (event.key === " ") {
          event.preventDefault(); // Prevent page scrolling
          api.togglePlayPause();
        }
      }
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [video]);

  const handleNext = () => {
    setOffset((prevOffset) => prevOffset + 1);
  };

  const handlePrevious = () => {
    setOffset((prevOffset) => prevOffset - 1);
  };

  const handleShowVideoList = () => {
    setShowVideoList(true);
    setTimeout(() => {
      videoListRef.current?.scrollIntoView({ behavior: "smooth" });
    }, 10);
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
      {info  && <p className="info">{info}</p>}
      {video && (
        <div className="video-container">
          <h2 className="video-title">{video.title}</h2>
          <div id="video-player"></div>
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
            <p className="hint">
            Press <b><i>F</i></b> to toggle fullscreen and <b><i>Space</i></b> to play/pause the video.
            </p>
        </div>
      )}

      <div className="pagination-buttons">
        <button onClick={handlePrevious} className="button">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
          Previous
        </button>
        <button onClick={handleNext} className="button">
          Next
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
        <button onClick={handleShowVideoList} className="button scroll-button">View All Videos</button>
      </div>

      {showVideoList && (
        <div ref={videoListRef} style={{ minHeight: "800px" }}>
          <VideoList onVideoClick={handleVideoClick} />
        </div>
      )}
    </div>
  );
};

export default App;
