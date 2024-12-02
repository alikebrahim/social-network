import React, { useContext } from "react";
import { Button } from "antd";
import { AuthContext } from "../AuthContext";

const ChatPage = () => {
  const { isLoggedIn } = useContext(AuthContext);

  if (!isLoggedIn) {
    return <div>You must be logged in to access this page.</div>;
  }

  const joinGroupChat = () => {
    alert("Joined the group chat!");
  };

  return (
    <div>
      <h2>Chat Page</h2>
      <Button type="primary" onClick={joinGroupChat}>
        Join Group Chat
      </Button>
    </div>
  );
};

export default ChatPage;
