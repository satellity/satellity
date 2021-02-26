## Satellity

Satellity is a 100% open source forum, written in Go. Please visit [https://routinost.com](https://routinost.com) for more details. For feedback, you can submit [issues](https://github.com/satellity/satellity/issues) or join our [slack](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)([https://bit.ly/31b6xeX](https://join.slack.com/t/satellity/shared_invite/enQtNTcwMTIyODAwMDgxLTNhYTUxMDgzMjE2NTcwYjRjMmE5ZDdjNzJjZTgzMjQ3MTQ2YzRiMDE3YTM4YmVjYzBjNTNjOGMxMmRhZTU2ZDM)), Let's learn Go together!


## How to deployment

1. A VPS, I'm using [Digital Ocean](https://m.do.co/c/f911b101f6ac) right now, and you can use any other VPS like GCP, AWS. You'll get some credit from the link, and it depends on you.
2. Install Nginx `sudo apt install nginx -y`, here is an example config of nginx https://github.com/satellity/satellity/blob/master/deploy/nginx_example.conf , I'm using Ubuntu 20.04 LTS
3. Install Postgresql `sudo apt install postgresql -y`, [how to install PostgreSQL On Ubuntu](https://www.digitalocean.com/community/tutorials/how-to-install-postgresql-on-ubuntu-20-04-quickstart), after create the database, you need import the database schema https://github.com/satellity/satellity/blob/master/internal/models/schema.sql
4. Deploy the api server and web, you can find the shell script here: https://github.com/satellity/satellity/tree/master/deploy
5. Use systemd to manage http server, here is the service template https://github.com/satellity/satellity/blob/master/deploy/systemd/satellity-http.service, and you can find the basic commands [here](https://wiki.archlinux.org/index.php/systemd)
6. Lets Encrypt, here is a step by step tutorial, https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-20-04
7. Then you can visit you website after restart nginx

Some people want to know why not use docker? The most direct reason is that I don't know docker and didn't take time on it. And a shell script is enough for me right now.


## Features

1. REST API back-end written in Golang
2. React-based frontend
3. PostgreSQL, one of the best open source, flexible database 
4. Social login (OAuth 2.0) only support Github now
5. JSON Web Tokens (JWT) are used for user authentication in the API
6. Markdown supported topic and comment
7. Model tested


## Built With

1. go version go1.15 darwin/amd64
2. postgres (PostgreSQL) 12.3
3. react ^16.13.1

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

## License

![https://opensource.org/licenses/MIT](https://img.shields.io/github/license/mashape/apistatus.svg)
