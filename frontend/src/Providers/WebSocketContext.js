import React, { createContext, useEffect, useState } from "react";

export const WebSocketContext = createContext(null);

let socket;

export const WebSocketProvider = ({ children }) => {
  const [ws, setWs] = useState(null);

  useEffect(() => {
    if (!socket) {
      socket = new WebSocket("ws://localhost:8080/ws");
      socket.onopen = () => {
        console.log("WebSocket connection established");
      };
      socket.onclose = () => {
        console.log("WebSocket connection closed");
      };
      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
      socket.onmessage = (message) => {
        console.log("WebSocket message received:", message.data);
      };
    }
    setWs(socket);
  }, []);

  return (
    <WebSocketContext.Provider value={ws}>{children}</WebSocketContext.Provider>
  );
};
