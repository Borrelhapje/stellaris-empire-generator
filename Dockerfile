FROM golang:1.19
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN GOPATH= GOARCH=wasm GOOS=js go build -o web/app.wasm
RUN GOPATH= go build .
FROM debian
COPY --from=0 /go/web/app.wasm /web/app.wasm
COPY --from=0 /go/stellaris .
CMD /stellaris