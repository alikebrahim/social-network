import React, { useContext } from "react";
import { Link } from "react-router-dom";
import { Menu } from "antd";
import { AuthContext } from "../AuthContext";

const Navbar = () => {
  const { isLoggedIn, logout } = useContext(AuthContext);

  const menuItems = [
    { key: "home", label: <Link to="/">Home</Link> },
    ...(isLoggedIn
      ? [
          { key: "chat", label: <Link to="/chat">Chat</Link> },
          { key: "add-post", label: <Link to="/add-post">Add Post</Link> },
          { key: "profile", label: <Link to="/profile">Profile</Link> },
          { key: "logout", label: <span onClick={logout}>Logout</span> },
        ]
      : [
          { key: "login", label: <Link to="/login">Login</Link> },
          { key: "signup", label: <Link to="/signup">Sign Up</Link> },
        ]),
  ];

  return <Menu mode="horizontal" theme="dark" items={menuItems} />;
};

export default Navbar;
