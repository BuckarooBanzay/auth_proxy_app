FROM ubuntu:artful

RUN apt-get update

# node stuff
RUN apt-get install -y nodejs npm

COPY package.json /data/
COPY src /data/src

RUN cd /data && npm i

WORKDIR /data

EXPOSE 8080

CMD ["node", "src/index.js"]
