FROM node:alpine

COPY package.json /data/
COPY package-lock.json /data/
COPY src /data/src

RUN cd /data && npm i

WORKDIR /data

EXPOSE 8080

CMD ["node", "src/index.js"]
