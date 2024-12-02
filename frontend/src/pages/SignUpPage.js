import React from "react";
import { Form, Input, Button } from "antd";

const SignUpPage = () => {
  const handleFinish = (values) => {
    console.log("User Signed Up:", values);
  };

  return (
    <Form onFinish={handleFinish} style={{ maxWidth: 400, margin: "auto" }}>
      <Form.Item
        name="email"
        rules={[{ required: true, message: "Please input your email!" }]}
      >
        <Input placeholder="Email" />
      </Form.Item>
      <Form.Item
        name="password"
        rules={[{ required: true, message: "Please input your password!" }]}
      >
        <Input.Password placeholder="Password" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Sign Up
        </Button>
      </Form.Item>
    </Form>
  );
};

export default SignUpPage;
