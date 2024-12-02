import React, { useContext } from "react";
import { Link } from "react-router-dom";
import { Menu } from "antd";
import { AuthContext } from "../AuthContext";

const Navbar = () => {
  const { isLoggedIn, logout } = useContext(AuthContext);

  return (
    <Menu mode="horizontal" theme="dark">
      <Menu.Item>
        <Link to="/">Home</Link>
      </Menu.Item>
      {isLoggedIn ? (
        <>
          <Menu.Item>
            <Link to="/chat">Chat</Link>
          </Menu.Item>
          <Menu.Item>
            <Link to="/add-post">Add Post</Link>
          </Menu.Item>
          <Menu.Item>
            <Link to="/profile">Profile</Link>
          </Menu.Item>
          <Menu.Item onClick={logout}>Logout</Menu.Item>
        </>
      ) : (
        <>
          <Menu.Item>
            <Link to="/login">Login</Link>
          </Menu.Item>
          <Menu.Item>
            <Link to="/signup">Sign Up</Link>
          </Menu.Item>
        </>
      )}
    </Menu>
  );
};

export default Navbar;
