import React from "react";
import { Card } from "antd";

const ProfilePage = () => {
  return (
    <Card title="Your Profile" style={{ maxWidth: 400, margin: "auto" }}>
      <p>Email: user@example.com</p>
      <p>Name: John Doe</p>
    </Card>
  );
};

export default ProfilePage;
