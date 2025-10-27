go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/go-redis/redis/v8
go get -u github.com/segmentio/kafka-go
go get -u github.com/joho/godotenv
go get -u github.com/golang-jwt/jwt/v5
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf

kafka-topics --create --topic order.created --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --create --topic order.status_changed --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --create --topic order.cancelled --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --list --bootstrap-server localhost:9092
## อ่าน messages
kafka-console-consumer --bootstrap-server localhost:9092 --topic order.created --from-beginning