# Social Network WebSocket Backend

## 📌 Overview

This is a real-time WebSocket API for a social media network built using Go. It allows users to **register, login, logout, create posts, follow/unfollow users, and retrieve followers**.

## 🚀 WebSocket Connection

To use the API, connect to the WebSocket server:

```
ws://localhost:8080/ws
```

## WebSocket Messages

### **1. Signup (Register New User)**

**Request:**

```json
{
  "type": "register",
  "data": {
    "email": "user@example.com",
    "password": "securepassword",
    "first_name": "John",
    "last_name": "Doe",
    "birth_date": "1990-01-01",
    "bio": "Hello, I'm John!",
    "avatar": "https://example.com/avatar.jpg",
    "nickname": "johndoe"
  }
}
```

**Expected Response:**

```json
{
  "type": "register_response",
  "content": "Signup successful"
}
```
- **If the email exists:**

```json
{
  "type": "register_response",
  "content": "Email already exists"
}
```


### **2. Login**

**Request:**

```json
{
  "type": "login",
  "data": {
    "email": "user@example.com",
    "password": "securepassword"
  }
}
```

**Expected Response:**

```json
{
  "type": "login_response",
  "content": "Login successful",
  "session_token": "your_session_token_here"
}
```

### **3. Logout**

**Request:**

```json
{
  "type": "logout",
  "data": {
    "session_token": "your_session_token_here"
  }
}
```

**Expected Response:**

```json
{
  "type": "logout_response",
  "content": "Logout successful"
}
```

- **If the session token is invalid:**

```json
{
  "type": "logout_response",
  "content": "Logout failed"
}
```

### **4. Create a Post**

**Request:**

```json
{
  "type": "create_post",
  "session_token": "your_session_token_here",
  "post": {
    "content": "Hello, this is my first post!",
    "image": "",
    "privacy": "public"
  }
}
```

**Expected Response:**

```json
{
  "type": "post_response",
  "content": "Post created successfully"
}
```

### ✅ **5. Follow Request**

**Request:**

```json
{
  "type": "follow_request",
  "session_token": "your_session_token_here",
  "follow": {
    "followed_id": 2
  }
}
```

**Expected Response:**

- **If the user has a public profile:**

```json
{
  "type": "follow_response",
  "content": "You are now following this user"
}
```

- **If the user has a private profile:**

```json
{
  "type": "follow_response",
  "content": "Follow request sent"
}
```

### ✅ **6. Get Followers**

**Request:**

```json
{
  "type": "get_followers",
  "session_token": "your_session_token_here"
}
```

**Expected Response:**

```json
{
  "type": "get_followers_response",
  "followers": [
    {
      "id": 1,
      "nickname": "john_doe",
      "avatar": "https://example.com/avatar1.jpg",
      "first_name": "John",
      "last_name": "Doe"
    },
    {
      "id": 2,
      "nickname": "jane_doe",
      "avatar": "https://example.com/avatar2.jpg",
      "first_name": "Jane",
      "last_name": "Doe"
    }
  ]
}
```

### ✅ **7. Unfollow a User**

**Request:**

```json
{
  "type": "unfollow",
  "session_token": "your_session_token_here",
  "follow": {
    "followed_id": 2
  }
}
```

**Expected Response:**

```json
{
  "type": "unfollow_response",
  "content": "Unfollowed successfully"
}
```

### ✅ **8. get post**

**Request:**

```json
{
  "type": "get_posts",
  "session_token": "your_valid_session_token"
}

```

**Expected Response:**

```json
{
  "type": "get_posts_response",
  "posts": [
    {
      "id": 1,
      "user_id": 2,
      "content": "This is a public post!",
      "image": "",
      "privacy": "public",
      "created_at": "2024-02-01 12:00:00"
    },
    {
      "id": 2,
      "user_id": 3,
      "content": "Private post, visible only to followers",
      "image": "",
      "privacy": "private",
      "created_at": "2024-02-01 12:05:00"
    }
  ]
}
```


## ⚡ **Next Steps**
 1. get post data 
 2. get user data 



