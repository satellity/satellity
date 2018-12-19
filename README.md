## Introduce

Godiscourse is a 100% open source, free, go back-end discourse like forum, build from `hello world`.

## Environment

1. GO version go1.10.2 darwin/amd64
2. postgres (PostgreSQL) 10.5

## Usage

1. `cd ./api`, copy `config/test.cfg` to `config/config.go`.
2. Prepare database, you can find database schema in `models/schema.sql`, find the database config in `config/config.cfg`.
3. Run back-end `go build` and `./api`.
4. All front-end is under `web` directory, `cd web` install dependence `npm install`, start service `npm run dev`.

## License

Released under the MIT license
