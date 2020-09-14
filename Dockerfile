# Stage 1 testing
FROM node:14.10.1-alpine

COPY .jshintrc /data/
COPY package.json /data/
COPY package-lock.json /data/
COPY src /data/src

RUN cd /data && npm i && npm test

# Stage 2 package
FROM node:14.10.1-alpine

COPY .jshintrc /data/
COPY package.json /data/
COPY package-lock.json /data/
COPY src /data/src

RUN cd /data && npm i --only=production

WORKDIR /data

EXPOSE 8080

CMD ["node", "src/index.js"]
