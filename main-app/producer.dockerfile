FROM gcr.io/distroless/nodejs AS runtime
EXPOSE 80
WORKDIR /tmp
COPY producer.js /
CMD ["/producer.js"]

