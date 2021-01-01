FROM node:14.5.0-alpine as builder

COPY . /data
WORKDIR /data

CMD ["sh", "-c", "npm i && npm test"]
