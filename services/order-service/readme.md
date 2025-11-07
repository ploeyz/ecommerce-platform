go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/go-redis/redis/v8
go get -u github.com/segmentio/kafka-go
go get -u github.com/joho/godotenv
go get -u github.com/golang-jwt/jwt/v5
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files


kafka-topics --create --topic order.created --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --create --topic order.status_changed --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --create --topic order.cancelled --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --list --bootstrap-server localhost:9092
## อ่าน messages
kafka-console-consumer --bootstrap-server localhost:9092 --topic order.created --from-beginning

Go plugins protpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user_service.proto