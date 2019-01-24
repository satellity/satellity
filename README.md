## Introduce

Godiscourse is a 100% open source, free, Golang back-end forum. The first code is `hello world`.

## Built With

1. Go version go1.10.2 darwin/amd64
2. postgres (PostgreSQL) 11.1
3. react ^16.7.0

## Structure

1. `api` is back-end service, which is Rails like structure.
2. `web` is front-end service, contains React, Parcel and etc.

## How to run it

1. `cd ./api`, copy `config/test.cfg` to `config/config.go`. Replace config with yours.
2. Prepare and start database, the database schema under `./api/models/schema.sql`, [how to install postgresql](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04).
3. `cd ./web` and run `npm install` to prepare front-end.
4. `cd path/to/api && go build && ./api` to start Golang server
5. `cd path/to/web && npm run dev` to run front-end, and open `http://localhost:1234`

## License

Released under the MIT license
