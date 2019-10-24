## Satellity

Satellity is a 100% open source forum, written in Go. Please visit [https://satellity.org](https://satellity.org) for more details. For feedback, you can submit [issues](https://github.com/satellity/satellity/issues) or join our [slack](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)([https://bit.ly/31b6xeX](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)), Let's learn Go together!

## NOTICE

Satellity is a still a **PRE-ALPHA** version. Please don't use it in production!!

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
2. `./app` is front-end service, contains React, Parcel and etc.
2. `./deploy` contains example of deploy, nginx and systemd.

## Screenshot

![Satellity](/screenshots/aspect.png "Hello Satellity")

## Getting Started

### Backend

1. `cd ./internal`, copy `config/config.example` to `config/config.yaml`. Replace config with yours.
2. Prepare and start database, the database schema under `./internal/models/schema.sql`, [how to install postgresql](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04).
3. `cd ./ && make install && make build && ./bin/satellity` to start Golang server

### Frontend

1. `cd ./app`, copy `.env.example` to `.env` and `.env.development` and change the following fields:
   
    ```
    SITE_NAME=your site name
    ```
2. run `yarn install` to prepare front-end.
3. `yarn start` and open `http://localhost:3000`

## Contribution
When contributing to this repository, please reach out to [@jadeydi](https://github.com/jadeydi) or other contributors via email, issue or any other means to discuss the changes you wish to make.

You can also just clone the repository, create a new branch of the feature or issue and make adequate changes then push and create a pull-request and request a review from other contributors.

## License

![https://opensource.org/licenses/MIT](https://img.shields.io/github/license/mashape/apistatus.svg)
