FROM node:20-alpine AS build

WORKDIR /visual-compiler-frontend

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /visual-compiler-frontend

COPY --from=build /visual-compiler-frontend/package.json ./
COPY --from=build /visual-compiler-frontend/.svelte-kit/output ./

EXPOSE 3000

CMD [ "node", "server/index.js" ]
