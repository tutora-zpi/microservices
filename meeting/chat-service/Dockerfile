FROM node:20-slim

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE 8002

CMD ["node", "dist/main"]
