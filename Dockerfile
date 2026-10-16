FROM golang:1.26.2-alpine as build

## Tạo folder để build
WORKDIR /app

## Copy src vào app
## Sử dụng chức năng lưu cache của docker 

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

## BUILD
RUN go build -o go-backend ./cmd

## alpine là hệ điều hành?
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/go-backend .

CMD ["./go-backend"]





## image lần đầu build. 2.26gb
## docker ignore ~ 2gb
## kỹ thuật stage  ~ 109MB GOOD

## sau khi build xong, chỉ cần file để chạy





