FROM node:18-alpine

WORKDIR /app

RUN apk add --no-cache python3 make g++

COPY package.json package-lock.json ./

RUN npm install
RUN npm rebuild bcrypt --build-from-source

COPY src ./src
COPY config ./config
COPY server.js .

ENV NODE_ENV=production
ENV CONFIG_PATH=/app/config/container.yaml

EXPOSE 5000

CMD ["npm", "start"]

