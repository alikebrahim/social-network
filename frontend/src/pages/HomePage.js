import React from "react";
import { Typography, Col, Row } from "antd";
import Post from "../components/Post";

const { Title } = Typography;

const HomePage = () => {
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
    // {
    //   id: 4,
    //   title: "Third Post",
    //   content: "This is the third post",
    //   image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
    //   likes: 20,
    //   dislikes: 0,
    //   comments: 8,
    // },
    // {
    //   id: 5,
    //   title: "Third Post",
    //   content:
    //     "This is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the This is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThisThis is the third postThis is the third postThisthird postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisThis is the third postThis is the third postThisis is the third postThisThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third postThis is the third post",
    //   // image: "https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png",
    //   likes: 20,
    //   dislikes: 0,
    //   comments: 8,
    // },
  ];

  return (
    <>
      <Title level={2}>Welcome to the Homepage!</Title>
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
    </>
  );
};

export default HomePage;
