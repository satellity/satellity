## api

api is the back-end of [GoDiscourse](https://live.satellity.com/), there isn't a convention for Go (maybe just i don't find), so I learned from RoR:

1. `controllers` contains routes of request.
2. `models` contains the methods of a struct, include read data from database.
3. `views` is where to place the response body, we only have JSON format here.
4. `middleware` is unique Go.
