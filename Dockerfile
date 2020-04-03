# Stage 1 testing
FROM node:13.12.0-alpine

COPY . /data
COPY .git/refs/heads/master /data/public/version.txt

RUN cd /data && npm i && npm test

# Stage 2 package
FROM node:13.12.0-alpine

COPY . /data
COPY .git/refs/heads/master /data/public/version.txt

RUN cd /data && npm i --only=production

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
