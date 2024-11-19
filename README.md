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

eg. `http://localhost:8080/posts/39`
