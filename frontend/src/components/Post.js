import React from "react";
import { Card, Button, Typography, Space, Tooltip } from "antd";
import {
  LikeOutlined,
  DislikeOutlined,
  CommentOutlined,
  EditOutlined,
} from "@ant-design/icons";

const { Text } = Typography;

const Post = ({ title, image, content, likes, dislikes, comments }) => {
  return (
    <Card
      title={title}
      hoverable
      style={{
        height: "500px", // Updated height
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
        // justifyContent: "space-between",
        margin: "10px",
        borderRadius: "10px",
        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
      }}
      cover={
        image && (
          <img
            src={image}
            alt={title}
            style={{
              height: "200px",
              objectFit: "contain",
              borderTopLeftRadius: "0",
              borderTopRightRadius: "0",
              borderBottomLeftRadius: "0",
              borderBottomRightRadius: "0",
            }}
          />
        )
      }
      actions={[
        <Space
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            padding: "10px",
          }}
        >
          <Space size="middle">
            <Tooltip title="Likes">
              <Space>
                <LikeOutlined />
                <Text>{likes}</Text>
              </Space>
            </Tooltip>
            <Tooltip title="Dislikes">
              <Space>
                <DislikeOutlined />
                <Text>{dislikes}</Text>
              </Space>
            </Tooltip>
            <Tooltip title="Comments">
              <Space>
                <CommentOutlined />
                <Text>{comments}</Text>
              </Space>
            </Tooltip>
          </Space>
          <Button type="primary" shape="round" icon={<EditOutlined />}>
            Edit
          </Button>
        </Space>,
      ]}
    >
      <Text>{content}</Text>
    </Card>
  );
};

export default Post;
