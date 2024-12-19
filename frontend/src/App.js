import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Layout } from "antd";
import { AuthProvider } from "./AuthContext";
import { WebSocketProvider } from "./Providers/WebSocketContext";
import Navbar from "./components/Navbar";
import HomePage from "./pages/HomePage";
import ChatPage from "./pages/ChatPage";
import LoginPage from "./pages/LoginPage";
import SignUpPage from "./pages/SignUpPage";
import AddPostPage from "./pages/AddPostPage";
import ViewPostPage from "./pages/ViewPostPage";
import ProfilePage from "./pages/ProfilePage";
import UpdateInfoPage from "./pages/UpdateInfoPage";
import UsersList from "./components/UsersList";

const App = () => {
  return (
    <AuthProvider>
      <WebSocketProvider>
        <Router>
          <Layout style={{}}>
            <Layout.Header
              style={{ position: "fixed", zIndex: 1, width: "100%" }}
            >
              <Navbar />
            </Layout.Header>
            <Layout.Content style={{ marginTop: 66 }}>
              <UsersList>
                <Routes>
                  <Route path="/" element={<HomePage />} />
                  <Route path="/login" element={<LoginPage />} />
                  <Route path="/signup" element={<SignUpPage />} />
                  <Route path="/chat/:type/:id" element={<ChatPage />} />
                  <Route path="/add-post" element={<AddPostPage />} />
                  <Route path="/view-post" element={<ViewPostPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="/update-info" element={<UpdateInfoPage />} />
                </Routes>
              </UsersList>
            </Layout.Content>
          </Layout>
        </Router>
      </WebSocketProvider>
    </AuthProvider>
  );
};

export default App;
