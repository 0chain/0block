module 0block/zblockcore

go 1.14

require (
	0block/core v0.0.0
	github.com/0chain/gosdk v1.1.4
	github.com/didip/tollbooth v4.0.2+incompatible // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace 0block/core => ../core
