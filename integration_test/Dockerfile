FROM node:15.5.1-alpine as builder

COPY . /data
WORKDIR /data

CMD ["sh", "-c", "npm i && npm test"]
