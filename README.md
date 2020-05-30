## Satellity

Satellity is a 100% open source forum, written in Go. Please visit [https://routinost.com](https://routinost.com) for more details. For feedback, you can submit [issues](https://github.com/satellity/satellity/issues) or join our [slack](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)([https://bit.ly/31b6xeX](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)), Let's learn Go together!

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

1. Go version go1.4 darwin/amd64
2. postgres (PostgreSQL) 11.4
3. react ^16.10.2

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
3. `cd ./ && go build && ./satellity` to start Golang server

### Frontend

1. Copy `env.example` to `.env`, and replace `Satellity` with your project name.
   
    ```
    SITE_NAME=your site name
    ```
2. run `yarn install`, then `yarn start`. It's running now.

## Contribution

When contributing to this repository, please reach out to [@jadeydi](https://github.com/jadeydi) or other contributors via email, issue or any other means to discuss the changes you wish to make.

You can also just clone the repository, create a new branch of the feature or issue and make adequate changes then push and create a pull-request and request a review from other contributors.

## Donation
If this project is helpful, you can also consider a small amount of donations.

1. [Paypal](https://www.paypal.me/jadeydi/5usd)
2. BTC Address: [1JXjQJ4tK7fsKf1biCisD4yKdm5PbWXkoD](https://imgur.com/a/TEtvQZ4)
3. ETH (or other erc20 token) Address: [0xAE9EA2D22E49B4c845Bbe57B57aB7172e548cE0B](https://imgur.com/a/hxM8YeF)

## License

![https://opensource.org/licenses/MIT](https://img.shields.io/github/license/mashape/apistatus.svg)
