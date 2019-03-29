# Build app within the base builder image
FROM node:8
LABEL maintainer="Guo Huang <guohuang@gmail.com>"

WORKDIR /src/godiscourse/web

# Copy solution
COPY package*.json ./

# Publish application for release
RUN npm install

COPY . .

EXPOSE 1234

CMD [ "npm", "run", "dev" ]