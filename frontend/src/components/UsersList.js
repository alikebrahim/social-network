import React, { useContext, useState } from "react";
import { Avatar, Input, Layout, List, Typography, Badge } from "antd";
import { AuthContext } from "../AuthContext";
import { useNavigate } from "react-router-dom";

const { Sider } = Layout;
const { Title } = Typography;

const UsersList = ({ children }) => {
  const { isLoggedIn } = useContext(AuthContext);
  const navigate = useNavigate();

  const [searchTerm, setSearchTerm] = useState("");
  const users = [
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
  ];

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
    <Layout style={{ height: "95vh" }}>
      <Sider width={300} style={{ background: "#fff", padding: "10px" }}>
        <Input
          placeholder="Search users and groups"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <Title level={3}>Users And Groups</Title>
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
                <div style={{ flex: 1, display: "flex", alignItems: "center" }}>
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
      </Sider>
      <Layout>{children}</Layout>
    </Layout>
  );
};

export default UsersList;
