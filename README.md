## CS567 Final Project:

### How to use
1. install [Docker](https://docs.docker.com/desktop/setup/install/mac-install/) and [docker-compose](https://docs.docker.com/compose/install/)
2. add .env file to current folder
   ```
   // .env
   MONGODB_URI=mongodb+srv://...
   ```
3. run `docker-compose up --build` to start the backend 

### API Documentation

This document provides details on the available API endpoints and their request/response formats.

---

### [GET] `/posts/{index}`
**Description:**  
Get the `{index}`th post from the database.

**Request:**
- Path Parameter:
  - `index`: The index of the post to retrieve.

**Response:**
- `200 OK`:
```json
{
  "_id": "string",
  "author": "string",
  "body": "string",
  "comments": [
    {
      "author": "string",
      "body": "string",
      "created_at": "string"
    }
  ],
  "created_at": "string",
  "id": "string",
  "permalink": "string",
  "title": "string",
  "upvote": 0,
  "url": "string"
}
```

---

### [GET] `/policy`
**Description:**  
Get the policy for this subreddit.

**Request:**
- None

**Response:**
- `200 OK`:
```json
{
  "policy_name": "string",
  "policy_description": "string",
  "post_id": "ObjectId",
  "vote_count": 0,
  "is_final": true
}
```

---

### [GET] `/newPolicy`
**Description:**  
Get a newly generated policy from the LLM.

**Request:**
- None

**Response:**
- `200 OK`:
```json
{
  "policy_name": "string",
  "policy_description": "string",
  "post_id": "ObjectId",
  "vote_count": 0,
  "is_final": false
}
```

---

### [POST] `/concern`
**Description:**  
Submit a concern for a post, and get the generated policy.

**Request:**
- Body (JSON):
```json
{
  "post": 42,
  "concern": "Data privacy issue"
}
```

**Response:**
- `200 OK`:
```json
{
  "policy": "Policy generated for post 000000000000000000000000 with concern: Data privacy issue"
}
```

---

### [POST] `/simulation`
**Description:**  
Submit a policy and its simulation data, and return the simulation results.

**Request:**
- Body (JSON):
```json
{
  "policy": "Some policy string",
  "simulation": "Some simulation string"
}
```

**Response:**
- `200 OK`:
```json
{
  "results": [
    {
      "role": "role string",
      "comment": "comment string"
    }
  ]
}
```

---

### [POST] `/vote`
**Description:**  
Submit a vote for a specific policy.

**Request:**
- Body (JSON):
```json
{
  "user": 1,
  "policy": "1qwewrw",
  "vote": "upvote"
}
```

**Response:**
- `200 OK`:
```json
{
  "success": true
}
```

---

### Endpoints Summary

- **GET /posts/{index}**  
  Get the `{index}`th post.

- **GET /policy**  
  Get policy for the subreddit.

- **GET /newPolicy**  
  Get new policy generated from the LLM.

- **POST /concern**  
  Submit a concern for a post and receive the generated policy.

- **POST /simulation**  
  Submit a policy and its simulation data, returning the results.

- **POST /vote**  
  Submit a vote for a policy.

  /// TODO
  // getPost: get all posts at a time
  // postComment: update comment of a post
  // postUpvote: update Upvote of a post

