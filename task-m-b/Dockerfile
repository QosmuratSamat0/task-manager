FROM node:18-alpine

WORKDIR /app

RUN apk add --no-cache python3 make g++

COPY package.json ./

RUN npm install
RUN npm rebuild bcrypt --build-from-source

COPY src ./src
COPY server.js .

ENV NODE_ENV=production

EXPOSE 5000

CMD ["npm", "start"]

