import React, { useContext, useState, useEffect } from "react";
import { Button, Input, Layout, Typography } from "antd";
import { AuthContext } from "../AuthContext";
import { useParams } from "react-router-dom";
import { users } from "../components/UsersList";
const { Header, Content } = Layout;
const { Title } = Typography;

const ChatPage = () => {
  const { type, id } = useParams();
  const { isLoggedIn } = useContext(AuthContext);
  const [selectedUser, setSelectedUser] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");

  if (!isLoggedIn) {
    return <div>You must be logged in to access this page.</div>;
  }

  const sendMessage = () => {
    if (newMessage.trim() !== "") {
      setMessages([...messages, { sender: "You", text: newMessage }]);
      setNewMessage("");
    }
  };

  useEffect(() => {
    const fetchInfo = () => {
      console.log(type, id);
      const user = users.find(
        (user) => user.id === parseInt(id) && user.type === type
      );
      setSelectedUser(user);
      setMessages([
        { sender: user.name, text: "newMessage" },
        { sender: user.name, text: "newMessage" },
        { sender: "You", text: "newMessage" },
        { sender: user.name, text: "newMessage" },
      ]);
    };
    fetchInfo();
  }, [id, type]);

  return (
    <Layout style={{ height: "95vh" }}>
      <Layout>
        <Header style={{ background: "#fff", padding: "0 20px" }}>
          <Title level={3}>
            {selectedUser
              ? `Chat with ${selectedUser.name}`
              : `Select a user to start chatting ${type} ${id}`}
          </Title>
        </Header>
        <Content
          style={{
            padding: "20px",
            display: "flex",
            flexDirection: "column",
            height: "100%",
          }}
        >
          {selectedUser && (
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                height: "100%",
              }}
            >
              <div
                style={{
                  border: "1px solid #ccc",
                  padding: "10px",
                  flex: 1,
                  overflowY: "auto",
                  marginBottom: "20px",
                  borderRadius: "5px",
                  backgroundColor: "#f9f9f9",
                }}
              >
                {messages.map((message, index) => (
                  <div
                    key={index}
                    style={{
                      marginBottom: "10px",
                      display: "flex",
                      justifyContent:
                        message.sender === "You" ? "flex-end" : "flex-start",
                    }}
                  >
                    <div
                      style={{
                        maxWidth: "60%",
                        padding: "10px",
                        borderRadius: "10px",
                        backgroundColor:
                          message.sender === "You" ? "#dcf8c6" : "#fff",
                        border:
                          message.sender === "You"
                            ? "1px solid #dcf8c6"
                            : "1px solid #ccc",
                      }}
                    >
                      <strong>{message.sender}:</strong> {message.text}
                    </div>
                  </div>
                ))}
              </div>
              <div style={{ display: "flex", alignItems: "center" }}>
                <Input
                  placeholder="Type a message"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  style={{ marginRight: "10px", flex: 1 }}
                />
                <Button type="primary" onClick={sendMessage}>
                  Send
                </Button>
              </div>
            </div>
          )}
        </Content>
      </Layout>
    </Layout>
  );
};

export default ChatPage;
