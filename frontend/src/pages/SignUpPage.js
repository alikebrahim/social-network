import React, { useContext, useState, useEffect } from "react";
import { Form, Input, Button, Alert, DatePicker } from "antd";
import { WebSocketContext } from "../Providers/WebSocketContext";
import { useNavigate } from "react-router-dom";
import sha256 from "crypto-js/sha256";

const SignUpPage = () => {
  const socket = useContext(WebSocketContext);
  const [error, setError] = useState(null);
  const [response, setResponse] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (socket) {
      socket.onmessage = (message) => {
        const data = JSON.parse(message.data);

        if (data["content"] === "Signup successful") {
          setError(null);
          setResponse(data["content"]);
          navigate("/login");
        } else {
          setResponse(null);
          setError(data["content"]);
        }
      };
    }
  }, [socket, navigate]);

  const handleFinish = (values) => {
    setResponse(null);
    if (
      values.first_name === "" ||
      values.last_name === "" ||
      values.email === "" ||
      values.password === "" ||
      values.birthday === null
    ) {
      setError("Please fill in all fields.");
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(values.email)) {
      setError("Please enter a valid email address.");
      return;
    }
    setError(null);
    try {
      if (socket) {
        values.password = sha256(values.password).toString();
        socket.send(JSON.stringify({ type: "register", data: values }));
      }
    } catch (err) {
      setError("An error occurred while signing up. Please try again.");
      console.error("Sign up error:", err);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: "auto" }}>
      {error && <br />}
      {response && <br />}
      {error && <Alert message={error} type="error" showIcon />}
      {response && (
        <Alert message={`Response: ${response}`} type="success" showIcon />
      )}
      <br />
      <Form onFinish={handleFinish}>
        <Form.Item
          name="first_name"
          rules={[{ required: true, message: "Please input your first name!" }]}
        >
          <Input placeholder="First Name" />
        </Form.Item>
        <Form.Item
          name="last_name"
          rules={[{ required: true, message: "Please input your last name!" }]}
        >
          <Input placeholder="Last Name" />
        </Form.Item>
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
        <Form.Item
          name="confirm_password"
          rules={[{ required: true, message: "Please input your password!" }]}
        >
          <Input.Password placeholder="Confirm Password" />
        </Form.Item>
        <Form.Item
          name="dob"
          rules={[
            { required: true, message: "Please input your date of birth!" },
          ]}
        >
          <DatePicker placeholder="Date of Birth" style={{ width: "100%" }} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Sign Up
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default SignUpPage;
