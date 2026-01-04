FROM node:20-alpine AS build-env
COPY . /app
WORKDIR /app
RUN corepack enable
RUN pnpm i
RUN pnpm build

FROM scratch
COPY --from=build-env /app/build/client /
