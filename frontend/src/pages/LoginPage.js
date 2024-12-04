import React, { useContext, useState, useEffect } from "react";
import { Form, Input, Button, Alert } from "antd";
import { AuthContext } from "../AuthContext";
import { WebSocketContext } from "../Providers/WebSocketContext";
import sha256 from "crypto-js/sha256";

const LoginPage = () => {
  const { login } = useContext(AuthContext);
  const socket = useContext(WebSocketContext);
  const [error, setError] = useState(null);
  const [response, setResponse] = useState(null);
  const [values, setValues] = useState({ email: "", password: "" });

  useEffect(() => {
    if (socket) {
      socket.onmessage = (message) => {
        const data = JSON.parse(message.data);
        console.log(data);
        if (data["content"] === "Login successful") {
          setError(null);
          setResponse(data["content"]);
          login({ email: values.email });
        } else {
          setResponse(null);
          setError(data["content"]);
        }
      };
    }
  }, [socket, login, values.email]);

  const handleFinish = (formValues) => {
    setValues(formValues);
    try {
      if (socket) {
        formValues.password = sha256(formValues.password).toString();
        socket.send(JSON.stringify({ type: "login", data: formValues }));
      }
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

  return (
    <div>
      {error && <br />}
      {response && <br />}
      {error && <Alert message={error} type="error" showIcon />}
      {response && (
        <Alert message={`Response: ${response}`} type="success" showIcon />
      )}
      <br />
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
    </div>
  );
};

export default LoginPage;
