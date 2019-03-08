## GoDiscourse
GoDiscourse is a 100% open source, discourse-like forum written in Go. For demo, please visit [https://live.godiscourse.com/](https://live.godiscourse.com/). For feedback, you can submit [issues](https://github.com/godiscourse/godiscourse/issues) or join our slack [https://bit.ly/2IV6LCW](https://join.slack.com/t/godiscourse/shared_invite/enQtNTcwMTIyODAwMDgxLWE4YjI4MWNiNGM0NDU5MGJiZTNjY2NiYzJhNjQ3NmUxZDA5NDU3ODg2NmY4ODM3NTcyZjIwYmM4OWFiZmEyNjE)

## Built With
1. Go version go1.11.4 darwin/amd64
2. postgres (PostgreSQL) 11.1
3. react ^16.7.0

## Features
1. REST API back-end written in Golang
2. React-based frontend
3. PostgreSQL, one of the best open source, flexible database 
4. Social login (OAuth 2.0) only support Github now
5. JSON Web Tokens (JWT) are used for user authentication in the API
6. Markdown supported topic and comment
7. Model tested

## Structure
1. `./api` is back-end service, which is Rails like structure.
2. `./web` is front-end service, contains React, Parcel and etc.
2. `./deploy` contains example of deploy, nginx and systemd.

## Getting Started
1. `cd ./api`, copy `config/test.cfg` to `config/config.go`. Replace config with yours.
2. Prepare and start database, the database schema under `./api/models/schema.sql`, [how to install postgresql](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04).
3. `cd ./web` and run `npm install` to prepare front-end.
4. `cd path/to/api && go build && ./api` to start Golang server
5. `cd path/to/web && npm run dev` to run front-end, and open `http://localhost:1234`

## License
![https://opensource.org/licenses/MIT](https://img.shields.io/github/license/mashape/apistatus.svg)
