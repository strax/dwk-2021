FROM node AS builder
WORKDIR /build
COPY package.json package-lock.json /build/
RUN npm ci
COPY . /build/
RUN npx webpack --env production

FROM nginx
COPY --from=builder /build/assets /usr/share/nginx/html/