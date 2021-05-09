FROM gcr.io/distroless/nodejs AS runtime
EXPOSE 80
WORKDIR /tmp
COPY consumer.js /
CMD ["/consumer.js"]

