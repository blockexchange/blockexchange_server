# Stage 1 testing
FROM node:15.5.1-alpine as builder

# deps
RUN apk add alpine-sdk cairo-dev pango-dev jpeg-dev

# files
COPY . /data

# compile / install
RUN cd /data &&\
  npm ci &&\
  npm test &&\
  npm run jshint_backend &&\
  npm run jshint_frontend &&\
  npm run bundle

# Stage 2 package
FROM node:15.5.1-alpine

COPY . /data

RUN apk add cairo pango jpeg
COPY --from=builder /data/node_modules /data/node_modules
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
