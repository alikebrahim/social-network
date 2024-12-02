import React from "react";
import { Form, Input, Button } from "antd";

const UpdateInfoPage = () => {
  const handleFinish = (values) => {
    console.log("Updated Info:", values);
  };

  return (
    <Form onFinish={handleFinish} style={{ maxWidth: 400, margin: "auto" }}>
      <Form.Item name="name">
        <Input placeholder="Name" />
      </Form.Item>
      <Form.Item name="email">
        <Input placeholder="Email" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update
        </Button>
      </Form.Item>
    </Form>
  );
};

export default UpdateInfoPage;
