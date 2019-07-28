## Satellity

Satellity is a 100% open source, written in Go. For forum part demo, please visit [https://live.godiscourse.com/](https://live.godiscourse.com/). For feedback, you can submit [issues](https://github.com/satellity/satellity/issues) or join our [slack](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLWE4YjI4MWNiNGM0NDU5MGJiZTNjY2NiYzJhNjQ3NmUxZDA5NDU3ODg2NmY4ODM3NTcyZjIwYmM4OWFiZmEyNjE)([https://bit.ly/2IV6LCW](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLWE4YjI4MWNiNGM0NDU5MGJiZTNjY2NiYzJhNjQ3NmUxZDA5NDU3ODg2NmY4ODM3NTcyZjIwYmM4OWFiZmEyNjE)), Let's learn Go together!

## Futrue

Satellity is a **PRE-ALPHA** version project. It will be a member management tool. Which contain two-part:

a. Private group, the main part of Satellity. which includes:

	1. Subscribe, can be paid by members, monthly, yearly, etc.
	2. Invitation, a member can be invited into group for free.
	3. Private messages are only visible for members.
	4. Hot messages, avoid messing anything important.
	5. Some group management functions still thinking
	
b. Forum, the public area anyone can share anything here.

It'll be a patreon-like forum.

## Features

1. REST API back-end written in Golang
2. React-based frontend
3. PostgreSQL, one of the best open source, flexible database 
4. Social login (OAuth 2.0) only support Github now
5. JSON Web Tokens (JWT) are used for user authentication in the API
6. Markdown supported topic and comment
7. Model tested


## Built With

1. Go version go1.12.5 darwin/amd64
2. postgres (PostgreSQL) 11.1
3. react ^16.8.4

## Structure

1. `./` is back-end service, we followed [golang-standards project-layout](https://github.com/golang-standards/project-layout).
2. `./web` is front-end service, contains React, Parcel and etc.
2. `./deploy` contains example of deploy, nginx and systemd.

## Screenshot

![Go Discourse](/screenshots/aspect.png "Hello Go Discourse")

## Getting Started

### Backend

1. `cd ./internal`, copy `config/config.yaml.example` to `config/config.yaml`. Replace config with yours.
2. Prepare and start database, the database schema under `./internal/models/schema.sql`, [how to install postgresql](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04).
3. `cd ./ && make install && make build && ./bin/satellity` to start Golang server

### Frontend

1. `cd ./web`, copy `.env.example` to `.env` and `.env.development` and change the following fields:
   
    ```
    SITE_NAME=your site name
    API_HOST=http://localhost:4000 or production url
    GITHUB_CLIENT_ID=put your client id
    ```
2. run `npm install` to prepare front-end.
3. `npm run dev` and open `http://localhost:1234`

## Docker

1. go to satellity folder
2. type `docker-compose up` 

## License

![https://opensource.org/licenses/MIT](https://img.shields.io/github/license/mashape/apistatus.svg)
