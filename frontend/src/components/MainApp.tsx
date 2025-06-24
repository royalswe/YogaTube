// MainApp.tsx
import React, { useState, useEffect } from "react";
import App from "./videos";
import "../main.css";

const MainApp: React.FC = () => {
  const [darkMode, setDarkMode] = useState(() => {
    const savedMode = localStorage.getItem("darkMode");
    return savedMode === "true";
  });

  useEffect(() => {
    document.body.className = darkMode ? "dark-mode" : "light-mode";
    localStorage.setItem("darkMode", darkMode.toString());
  }, [darkMode]);

  const toggleMode = () => {
    setDarkMode((prevMode) => !prevMode);
  };

  return (
    <div>
      <button onClick={toggleMode} className="toggle-button">
        {darkMode ? "Switch to Light Mode" : "Switch to Dark Mode"}
      </button>
      <App />
    </div>
  );
};

export default MainApp;
