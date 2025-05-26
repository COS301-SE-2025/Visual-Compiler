FROM node:20-alpine

WORKDIR /visual-compiler-frontend

COPY package.json ./
RUN npm install

COPY . .

EXPOSE 5173

CMD [ "npm", "run", "dev", "--", "--host" ]
