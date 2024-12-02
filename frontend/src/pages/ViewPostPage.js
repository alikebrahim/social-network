import React from "react";
import { List } from "antd";

const posts = [
  { title: "Post 1", content: "This is the first post" },
  { title: "Post 2", content: "This is the second post" },
];

const ViewPostPage = () => {
  return (
    <List
      dataSource={posts}
      renderItem={(item) => (
        <List.Item>
          <List.Item.Meta title={item.title} description={item.content} />
        </List.Item>
      )}
    />
  );
};

export default ViewPostPage;
