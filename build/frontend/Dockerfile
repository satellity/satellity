FROM node:lts-alpine as builder

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY app ./

RUN mv env.example .env \
    && yarn install

RUN yarn build

FROM nginx:alpine

COPY --from=builder /usr/src/app/build /usr/share/nginx/html
