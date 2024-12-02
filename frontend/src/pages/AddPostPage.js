import React from "react";
import { Form, Input, Button } from "antd";

const AddPostPage = () => {
  const handleFinish = (values) => {
    console.log("Post Added:", values);
  };

  return (
    <Form onFinish={handleFinish} style={{ maxWidth: 600, margin: "auto" }}>
      <Form.Item
        name="title"
        rules={[{ required: true, message: "Please input the title!" }]}
      >
        <Input placeholder="Title" />
      </Form.Item>
      <Form.Item
        name="content"
        rules={[{ required: true, message: "Please input the content!" }]}
      >
        <Input.TextArea rows={4} placeholder="Content" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Add Post
        </Button>
      </Form.Item>
    </Form>
  );
};

export default AddPostPage;
