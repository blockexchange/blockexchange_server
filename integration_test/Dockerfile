FROM node:15.5.0-alpine as builder

# files
COPY . /data

# compile / install
RUN cd /data && npm ci

WORKDIR /data

EXPOSE 8080

CMD ["npm", "test"]
