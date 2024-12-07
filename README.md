## CS567 Final Project:

### How to use
1. install [Docker](https://docs.docker.com/desktop/setup/install/mac-install/) and [docker-compose](https://docs.docker.com/compose/install/)
2. add .env file to current folder
   ```
   // .env
   MONGO_URI=mongodb+srv://...
   GEMINI_API_KEY=...
   ```
3. run `docker-compose up --build` to start the backend 

### API Documentation

This document provides details on the available API endpoints and their request/response formats.

---

### [GET] `/posts`
**Description:**  
Get all posts from the database.

**Request:**
- None

**Response:**
- `200 OK`:
```json
[
  {
    "_id": "673bd2a759737946ba94048c",
    "author": "peachyisonline",
    "body": "Currently taking an 8-week class, and things from the first week still haven't been graded. Not quite sure how to address this, since it's hard to see if I'm on the right path with my assignments if I never receive any feedback. What's the appropriate resource to reach out to?",
    "comments": [
      {
        "author": "Budget-Reception-882",
        "body": "Email your professor explaining that’s it’s near the end of the semester and you’d like to know where you stand in their course.",
        "created_at": "2024-11-17T18:49:37"
      },
      {
        "author": "princam_",
        "body": "You can email them and risk annoying them, or just accept it because apparently that is the standard here.",
        "created_at": "2024-11-17T19:17:21"
      }
    ],
    "created_at": "2024-11-17T18:32:03",
    "id": "t3_1gts5ny",
    "permalink": "/r/UIUC/comments/1gts5ny/what_to_do_about_lack_of_grading/",
    "title": "What to do about lack of grading",
    "upvote": 9,
    "url": "https://www.reddit.com/r/UIUC/comments/1gts5ny/what_to_do_about_lack_of_grading/"
  },
  ...
]
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
[
  {
    "_id": "673e70dbfb491a41cf2c3680",
    "policy_name": "Purpose of the subreddit",
    "policy_description": "This subreddit is primarily for any topic of interest to students, staff, faculty, alumni, etc of UIUC. Additionally, anything related to Urbana or Champaign is welcome.",
    "vote_count": 0,
    "is_final": true
  },
  ...
]
```

---

### [GET] `/newPolicy`
**Description:**  
Get all newly generated policy from the LLM.

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
    "userID": "tester3",
    "_id":"673bd2a759737946ba94049d",
    "concern": "shouldn't ban these post"
}
```

**Response:**
- `200 OK`:
```json
{
    "policy": "Posts encouraging or glorifying harmful behavior, such as excessive drinking, will be removed.  This includes posts that solicit comparisons of negative experiences related to such behavior.",
    "simulateResult": {
        "Abuser": "The abuser would craft a post subtly comparing negative consequences of *responsible* alcohol consumption (e.g., a mild hangover) to highlight the perceived absurdity of the rule, thereby indirectly mocking the policy's overreach without explicitly glorifying excessive drinking.  This approach aims to exploit a loophole by focusing on the unintended consequences of the policy's broad interpretation.\n",
        "Moderator": "The moderator will proactively monitor posts for content glorifying or encouraging excessive drinking, including requests for negative experience comparisons, and promptly remove any violating submissions with a clear explanation of the policy violation.  Additionally, the moderator will educate the community about the new policy through announcements and responses to relevant queries, emphasizing the importance of responsible content creation.\n",
        "Policy Overview": "The policy, while well-intentioned, risks inconsistent enforcement due to vague terminology, potentially leading to frustration and circumvention.  Abusers may exploit ambiguities to undermine its purpose, creating a need for clearer definitions and more specific examples.  Proactive moderation is crucial, but educational efforts must emphasize the spirit, not just the letter, of the policy.  Improved clarity and community engagement are key to successful implementation.\n",
        "Regular User": "The general user finds this policy reasonable as it discourages potentially harmful behavior; however,  the vagueness of \"excessive drinking\" and \"negative experiences\" might lead to inconsistent enforcement.\n"
    },
    "success": true
}
```

---

### [GET] `/simulation`
**Description:**  
Submit a policy and its simulation data, and return the simulation results.

**Request:**
- none

**Response:**
- `200 OK`:
```json
[
    {
        "policy": "Posts encouraging or glorifying harmful behavior, including substance abuse, are prohibited.  Posts soliciting comparisons of personal struggles related to substance abuse are also prohibited.",
        "result": {
            "Abuser": "The abuser would subtly ask for advice on managing stress, framing it as a coping mechanism without explicitly mentioning substance abuse, then subtly weave in details about unhealthy coping methods to gauge responses and indirectly solicit comparisons.  This allows for a back-and-forth discussion that hints at substance abuse without directly glorifying or explicitly soliciting comparisons, relying on community members to fill in the gaps.\n",
            "Moderator": "The moderator will clearly announce the new policy to the subreddit, highlighting examples of prohibited content and emphasizing the consequences of violations.  Then, the moderator will actively monitor posts and comments, removing violations and issuing warnings or bans as needed, while also engaging with the community to foster understanding and compliance.\n",
            "Policy Overview": "The new policy effectively protects vulnerable users but risks chilling legitimate discussions about substance abuse support.  Clearer guidelines distinguishing between harmful glorification and requests for help are crucial.  Moderator training on nuanced identification of veiled requests for assistance is essential for policy effectiveness.  Proactive community engagement, perhaps through informative posts and FAQs, can mitigate unintended consequences.\n",
            "Regular User": "The general user finds this policy reasonable as it discourages potentially harmful behavior and protects vulnerable individuals, but worries it might unintentionally stifle discussions about seeking help for substance abuse.\n"
        }
    }
]
```

---

### [POST] `/updatePolicy`
**Description:**  
Submit a vote/comment for a specific policy.

**Request:**
- Body (JSON):
```json
{
    "_id": "674e1ce7c0992167e68d6601",
    "userID": "tester2",
    "vote": 1,
    "comment": "this is great."
}
```

**Response:**
- `200 OK`:
```json
{
    "message": "Policy updated successfully",
    "success": true
}
```

---

### [POST] `/updatePost`
**Description:**  
Update contents of a post.

**Request**
- Body (JSON):
```json
{
    "userID": "tester1",
    "post": {
    "_id": "673bd2a759737946ba94048c",
    "author": "MODIFIEDLOL",
    "body": "Currently taking an 8-week class, and things from the first week still haven't been graded. Not quite sure how to address this, since it's hard to see if I'm on the right path with my assignments if I never receive any feedback. What's the appropriate resource to reach out to?",
    "comments": [
      {
        "author": "Budget-Reception-882",
        "body": "Email your professor explaining that’s it’s near the end of the semester and you’d like to know where you stand in their course.",
        "created_at": "2024-11-17T18:49:37"
      },
      {
        "author": "princam_",
        "body": "You can email them and risk annoying them, or just accept it because apparently that is the standard here.",
        "created_at": "2024-11-17T19:17:21"
      }
    ],
    "created_at": "2024-11-17T18:32:03",
    "id": "t3_1gts5ny",
    "permalink": "/r/UIUC/comments/1gts5ny/what_to_do_about_lack_of_grading/",
    "title": "What to do about lack of grading",
    "upvote": 9,
    "url": "https://www.reddit.com/r/UIUC/comments/1gts5ny/what_to_do_about_lack_of_grading/"
  }
}
```

**Response:**
- `200 OK`:
```json
{
    "message": "Post updated successfully",
    "success": true
}
```

### [GET] `/log`
**Description:**  
Get log.

**Request:**
- None

**Response:**
- `200 OK`:
```json
{
    "userID": "tester_no",
    "timestamp": "2024-12-06T23:45:02.348Z",
    "level": "POSTCONCERN",
    "message": "Processed post concern request",
    "request": {
        "_id": "673bd2a759737946ba94048c",
        "concern": "should ban these nonsense post",
        "userID": "tester_no"
    },
    "response": {
        "policy": "",
        "success": false
    }
}
```

### Endpoints Summary

- **GET /posts**  
  Get all posts.

- **GET /policy**  
  Get policy for the subreddit.

- **GET /newPolicy**  
  Get new policy generated from the LLM.

- **GET /simulation**  
  Get the simulation results of new policies.

- **POST /concern**  
  Submit a concern for a post, simulate and output the generated policy. If no need to produce new policy or produced new policy doesn't pass simulation test, no policy will return.

- **POST /updatePolicy**  
  Submit a vote or comment a policy.

- **POST /updatePost**  
  Submit a update of a post.
