## CS567 Final Project:

### How to use
1. install [Docker](https://docs.docker.com/desktop/setup/install/mac-install/) and [docker-compose](https://docs.docker.com/compose/install/)
2. add .env file to current folder
   ```
   // .env
   MONGODB_URI=mongodb+srv://...
   ```
3. run `docker-compose up` to start the backend 

### API
#### [GET] /posts/{index}
Get {index}th post.
#### [GET] /policy Get policy for this subreddit. (JSON)
#### [GET] /newPolicy Get new policy generated from our LLM (JSON)
#### [POST] (JSON {“post”: Int, “concern”: String}, {“policy”: String})
#### [POST] (JSON {“policy”: String, “Simulation”:{}}, {“results”:[{“role”:String,’’comment”: String}]})
#### [POST] /vote (JSON {“user”: Int, “policy”: Int, “vote”: String})



eg. `http://localhost:8080/posts/39`
