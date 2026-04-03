FROM debian:bookworm-slim
WORKDIR /app
COPY src/daily-dish-server .
RUN chmod +x daily-dish-server
ENV PORT=8000
EXPOSE 8000
CMD ["./daily-dish-server"]
