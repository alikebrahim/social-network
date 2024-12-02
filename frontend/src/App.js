import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./AuthContext";
import { WebSocketProvider } from "./Providers/WebSocketContext";
import Navbar from "./components/Navbar";
import Notification from "./components/Notification";
import HomePage from "./pages/HomePage";
import ChatPage from "./pages/ChatPage";
import LoginPage from "./pages/LoginPage";
import SignUpPage from "./pages/SignUpPage";
import AddPostPage from "./pages/AddPostPage";
import ViewPostPage from "./pages/ViewPostPage";
import ProfilePage from "./pages/ProfilePage";
import UpdateInfoPage from "./pages/UpdateInfoPage";

const App = () => {
  return (
    <AuthProvider>
      <WebSocketProvider>
        <Router>
          <Navbar />
          <Notification />
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/signup" element={<SignUpPage />} />
            <Route path="/chat" element={<ChatPage />} />
            <Route path="/add-post" element={<AddPostPage />} />
            <Route path="/view-post" element={<ViewPostPage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/update-info" element={<UpdateInfoPage />} />
          </Routes>
        </Router>
      </WebSocketProvider>
    </AuthProvider>
  );
};

export default App;
