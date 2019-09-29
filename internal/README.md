## api

api is the back-end of [Satellity](https://live.godiscourse.com/), there isn't a convention for Go (maybe just I don't find), so I learned from RoR:

1. `controllers` contains the routers of request.
2. `models` contains the main methods, include operate data from database (CRUD).
3. `views` is where to place the response body, we only have JSON format here.
4. `middlewares` is unique Go.
