FROM node:carbon-alpine

# Create app directory
WORKDIR /usr/src/app

RUN mkdir /usr/src/pb

COPY pb/. /usr/src/pb/

# Install app dependencies
COPY mqtt-server/package.json mqtt-server/yarn.lock ./

RUN apk --no-cache add --virtual native-deps \
  g++ gcc libgcc libstdc++ linux-headers make python && \
  npm install --quiet node-gyp -g && \
  yarn install --production && \
  apk del native-deps

# Bundle app source
COPY mqtt-server/. .

EXPOSE 1883
CMD [ "npm", "start" ]
