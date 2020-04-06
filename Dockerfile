# Stage 1 testing
FROM node:13.12.0-alpine

COPY . /data
COPY .git/refs/heads/master /data/public/version.txt

# build
RUN cd /data &&\
  npm i &&\
  npm run jshint_backend &&\
  npm run jshint_frontend &&\
  npm run bundle


# Stage 2 package
FROM node:13.12.0-alpine

COPY . /data
COPY .git/refs/heads/master /data/public/version.txt

RUN cd /data && npm i --only=production
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
