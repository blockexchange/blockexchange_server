FROM node:15.7.0-alpine as builder

COPY . /data
WORKDIR /data

CMD ["sh", "-c", "npm i && npm test"]
