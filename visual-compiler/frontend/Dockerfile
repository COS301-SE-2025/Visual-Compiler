FROM node:20-alpine AS build

WORKDIR /visual-compiler-frontend

COPY package.json package-lock.json* ./
RUN npm install --frozen-lockfile

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /visual-compiler-frontend

COPY --from=build /visual-compiler-frontend/package.json ./
COPY --from=build /visual-compiler-frontend/build ./build

RUN npm install --omit=dev

EXPOSE 3000

CMD [ "node", "build" ]
