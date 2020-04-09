# Stage 1 testing
FROM node:13.12.0-alpine as builder

COPY . /data

# build
RUN cd /data &&\
  npm i &&\
  npm test &&\
  npm run jshint_backend &&\
  npm run jshint_frontend &&\
  npm run bundle

# use optimized js bundle
RUN cp /data/public/index_prod.html /data/public/index.html

# Stage 2 package
FROM node:13.12.0-alpine

COPY . /data
RUN apk update && apk add curl

RUN cd /data && npm i --only=production
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1

CMD ["npm", "start"]
