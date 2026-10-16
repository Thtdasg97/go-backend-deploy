## Docker build
Check tiền đề thì có dockerfile


docker build -t img-go-backend:latest .

## Docker clear cache
docker builder prune -af

## Giờ mới biết là Docker có lưu cache nữa đó :)

Dùng dive để xem image là gì :) 
dùng để debug image

## Kỹ thuật Stage - tối ưu dung lượng?

## Chạy img
docker run --name con-go-backend -d -p 3066:3069 --env-file .env.production img-go-backend:latest

## xem log

docker logs -f con-go-backend


docker container remove con-go-backend
docker container stop con-go-backend
docker container start con-go-backend
docker container restart con-go-backend

## Docker compose
Không lẽ cứ mỗi lần làm là phải đi tìm mấy cái ip bằng tay thế à? No, vì ta có docker compose
Nó sẽ không tìm bằng IP nữa, mà tìm bằng tên của service hay nói cách dễ hiểu, nó sẽ gọi thông qua tên service thay vì gọi thông qua IP

## Run docker compose
docker compose --env-file .env.production up -d

## tạo network
docker network create go-network

## down
docker compose down


## lấy authen của elastic search
docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana


## lấy finger print es
docker exec go-elastic cat /usr/share/elasticsearch/config/certs/http_ca.crt | openssl x509 -fingerprint -sha256 -noout


## khi xóa, thêm -v nghĩa là xóa luôn volume


## Câu lệnh tạo lại container đối tượng mà không xóa các container ngoài đối tượng
docker compose --env-file .env.production up -d --force-recreate go-backend
