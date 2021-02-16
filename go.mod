module app2

go 1.15

replace github.com/khunmk/pubsub => ./pubsub

require (
	github.com/gorilla/websocket v1.4.2
	github.com/khunmk/pubsub v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
)
