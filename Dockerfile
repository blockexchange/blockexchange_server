# Stage 1 testing
FROM node:13.10.1-alpine

COPY package.json /data/
COPY package-lock.json /data/
COPY src /data/src
COPY public /data/public
COPY .git/refs/heads/master /data/public/version.txt

RUN cd /data && npm i && npm test

# Stage 2 package
FROM node:13.10.1-alpine

COPY package.json /data/
COPY package-lock.json /data/
COPY src /data/src
COPY public /data/public
COPY .git/refs/heads/master /data/public/version.txt

RUN cd /data && npm i --only=production

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
