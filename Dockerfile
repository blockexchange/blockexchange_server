# Stage 1 testing
FROM node:15.4.0-alpine as builder

COPY . /data

# build
RUN apk add alpine-sdk cairo-dev pango-dev jpeg-dev &&\
	cd /data &&\
  npm ci &&\
  npm test &&\
  npm run jshint_backend &&\
  npm run jshint_frontend &&\
  npm run bundle

# Stage 2 package
FROM node:15.4.0-alpine

COPY . /data

RUN apk add cairo pango jpeg
COPY --from=builder /data/node_modules /data/node_modules
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
