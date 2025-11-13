FROM node:18-alpine

WORKDIR /app

COPY server.js package.json ./

RUN npm install

EXPOSE 8080

CMD ["node", "server.js"]
