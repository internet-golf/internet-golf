# i'll be real, i have no idea why you would need this many stages, but this is
# what the official dockerfile looks like

FROM node:20-alpine AS development-dependencies-env
COPY . /app
WORKDIR /app
RUN corepack enable
RUN pnpm i

FROM node:20-alpine AS production-dependencies-env
COPY ./package.json ./pnpm-lock.yaml /app/
WORKDIR /app
RUN npm ci --omit=dev

FROM node:20-alpine AS build-env
COPY . /app/
COPY --from=development-dependencies-env /app/node_modules /app/node_modules
WORKDIR /app
RUN npm run build

FROM scratch
COPY --from=build-env /app/build/client /
