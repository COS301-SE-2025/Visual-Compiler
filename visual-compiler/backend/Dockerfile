FROM alpine:latest

WORKDIR /visual-compiler-backend

COPY backend .

RUN chmod +x ./backend

ENV GOGC=50
ENV GOMAXPROCS=1

EXPOSE 8080

CMD [ "./backend" ]