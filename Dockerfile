# 배포용 컨테이너에 포함시킬 바이너리를 생성하는 컨테이너
FROM golang:1.22-alpine as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# ----------------------
# 배포용 컨테이너
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# ----------------------

# 로컬 개발 환경에서 사용하는 자동 새로고침 환경
FROM golang:1.22-alpine as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]