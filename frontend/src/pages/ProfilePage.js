import React from "react";
import { Card, Button, Row, Col, Typography } from "antd";
import Post from "../components/Post";
const { Title, Text } = Typography;

const ProfilePage = () => {
  const user = {
    id: 1,
    name: "Alice",
    followers: 10,
    following: 50,
  };

  const posts = [
    {
      id: 1,
      title: "Welcome",
      content: "This is the first post",
      image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
      likes: 10,
      dislikes: 2,
      comments: 5,
    },
    {
      id: 2,
      title: "Second Post",
      content: "This is the second post",
      // image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
      likes: 15,
      dislikes: 1,
      comments: 3,
    },
    {
      id: 3,
      title: "Third Post",
      content: "This is the third post",
      // image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
      likes: 20,
      dislikes: 0,
      comments: 8,
    },
    {
      id: 4,
      title: "Third Post",
      content: "This is the third post",
      image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
      likes: 20,
      dislikes: 0,
      comments: 8,
    },
    {
      id: 5,
      title: "Third Post",
      content:
        "This is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the This is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisthird postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisis is the third postThisThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third post",
      // image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
      likes: 20,
      dislikes: 0,
      comments: 8,
    },
  ];

  return (
    <div style={{ padding: "20px" }}>
      <div
        style={{ display: "flex", alignItems: "center", marginBottom: "20px" }}
      >
        <div style={{ marginLeft: "20px" }}>
          <Title level={2}>{user.name}</Title>
          <div style={{ display: "flex", gap: "20px" }}>
            <Text strong style={{ fontSize: "18px" }}>
              Followers: {user.followers}
            </Text>
            <Text strong style={{ fontSize: "18px" }}>
              Following: {user.following}
            </Text>
          </div>
        </div>
        <Button type="primary" style={{ marginLeft: "auto" }}>
          Edit
        </Button>
      </div>
      <Title level={3}>Posts</Title>
      <Row gutter={[16, 16]}>
        {posts.map((post) => (
          <Col key={"post.id"} xs={24} sm={12} md={8}>
            <Post
              title={post.title}
              comments={post.comments}
              likes={post.likes}
              image={post.image}
              dislikes={post.dislikes}
              content={post.content}
            ></Post>
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default ProfilePage;
