import React, { useContext } from "react";
import { Form, Input, Button } from "antd";
import { AuthContext } from "../AuthContext";

const LoginPage = () => {
  const { login } = useContext(AuthContext);

  const handleFinish = (values) => {
    login({ email: values.email });
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
          Login
        </Button>
      </Form.Item>
    </Form>
  );
};

export default LoginPage;
