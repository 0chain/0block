module zblock

go 1.14

require (
	0block/core v0.0.0
	0block/zblockcore v0.0.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/spf13/viper v1.6.3
	go.uber.org/zap v1.15.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace 0block/core => ../../core

replace 0block/zblockcore => ../../zblockcore
