## API

API is the back-end of [Routinost](https://routinost.com), there isn't a convention for Go (maybe just I don't find), so I learned from RoR:

1. `clouds` includes third part services, like mailgun and recaptcha.
2. `configs` contains server side configs.
3. `controllers` contains the routers of requests.
4. `durable` contains database.
5. `middlewares` is unique Go.
6. `models` contains the functions, include operate data from database (CRUD).
7. `session` contains all errors.
8. `views` is where place the response body, we only have JSON format here.
