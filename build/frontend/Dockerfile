FROM node:lts-alpine as builder

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY app ./

RUN mv env.example .env \
    && yarn install

RUN yarn build

COPY deploy/nginx_example.conf /usr/src/conf/nginx.conf

RUN sed -i 's#/path/to/satellity/html#/usr/share/nginx/html#g' /usr/src/conf/nginx.conf

FROM nginx:alpine

COPY --from=builder /usr/src/app/build /usr/share/nginx/html
COPY --from=builder /usr/src/conf/nginx.conf /etc/nginx/nginx.conf
