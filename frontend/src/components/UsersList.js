import React, { useContext, useState } from "react";
import { Avatar, Input, Layout, List, Typography, Badge } from "antd";
import { AuthContext } from "../AuthContext";
import { useNavigate } from "react-router-dom";

const { Sider } = Layout;
const { Title } = Typography;

export const users = [
  { id: 1, name: "Alice", isOnline: true, notifications: 2, type: "user" },
  { id: 2, name: "Bob", isOnline: false, notifications: 0, type: "user" },
  { id: 3, name: "Charlie", isOnline: true, notifications: 5, type: "user" },
  {
    id: 4,
    name: "adsadas 1",
    isOnline: true,
    notifications: 3,
    type: "group",
  },
  {
    id: 5,
    name: "asdasdas 2",
    isOnline: true,
    notifications: 0,
    type: "group",
  },
  { id: 2, name: "Bob1", isOnline: false, notifications: 0, type: "user" },
  { id: 2, name: "Bob12", isOnline: false, notifications: 0, type: "user" },
  { id: 2, name: "Bob123", isOnline: false, notifications: 0, type: "user" },
  { id: 2, name: "Bob1234", isOnline: false, notifications: 0, type: "user" },
  { id: 2, name: "Bob12345", isOnline: false, notifications: 0, type: "user" },
  { id: 2, name: "Bob123456", isOnline: false, notifications: 0, type: "user" },
  {
    id: 2,
    name: "Bob1234567",
    isOnline: false,
    notifications: 0,
    type: "user",
  },
  {
    id: 2,
    name: "Bob12345678",
    isOnline: false,
    notifications: 0,
    type: "user",
  },
  {
    id: 2,
    name: "Bob123456789",
    isOnline: false,
    notifications: 0,
    type: "user",
  },
  {
    id: 2,
    name: "Bob1234567890",
    isOnline: false,
    notifications: 0,
    type: "user",
  },
];

const UsersList = ({ children }) => {
  const { isLoggedIn } = useContext(AuthContext);
  const navigate = useNavigate();

  const [searchTerm, setSearchTerm] = useState("");

  const selectUser = (user) => {
    navigate(`/chat/${user.type}/${user.id}`);
  };

  const filteredUsers = users.filter((user) =>
    user.name.toLowerCase().includes(searchTerm.toLowerCase())
  );
  if (!isLoggedIn) {
    return <div>{children}</div>;
  }

  return (
    <Layout style={{}}>
      <Sider
        width={300}
        style={{
          background: "#fff",
          padding: "10px",
          overflow: "hidden",
        }}
      >
        <Input
          placeholder="Search users and groups"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <Title level={3}>Users And Groups</Title>
        <div
          style={{
            overflowY: "auto",
            height: "calc(100vh - 186px)",
          }}
        >
          <List
            dataSource={filteredUsers}
            renderItem={(user) => (
              <List.Item key={user.id} onClick={() => selectUser(user)}>
                <div
                  style={{
                    cursor: "pointer",
                    marginBottom: "10px",
                    display: "flex",
                    alignItems: "center",
                    width: "100%",
                    padding: "10px",
                    backgroundColor: "#f0f0f0",
                    borderRadius: "5px",
                    transition: "background-color 0.3s",
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = "#e6f7ff";
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = "#f0f0f0";
                  }}
                >
                  <div
                    style={{ flex: 1, display: "flex", alignItems: "center" }}
                  >
                    <Avatar
                      style={{
                        backgroundColor:
                          user.type === "group"
                            ? "#000000"
                            : user.isOnline
                            ? "#52c41a"
                            : "#f5222d",
                        marginRight: "10px",
                      }}
                    >
                      {user.name.charAt(0)}
                    </Avatar>
                    <span>
                      {user.name}
                      <br />
                      <span style={{ fontSize: "10px", color: "#888" }}>
                        {user.type === "group"
                          ? "GROUP"
                          : user.isOnline
                          ? "ONLINE"
                          : "OFFLINE"}
                      </span>
                    </span>
                  </div>
                  <div style={{ display: "flex", alignItems: "center" }}>
                    {user.notifications > 0 && (
                      <Badge
                        count={user.notifications}
                        style={{
                          backgroundColor: "#f5222d",
                          marginLeft: "10px",
                        }}
                      />
                    )}
                  </div>
                </div>
              </List.Item>
            )}
          />
        </div>
      </Sider>
      <Layout>{children}</Layout>
    </Layout>
  );
};

export default UsersList;
