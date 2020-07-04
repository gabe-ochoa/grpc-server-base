FROM golang:1.14-alpine3.12

COPY ./dist /dist

CMD [ "/dist/grpc-server-api-linux" ]