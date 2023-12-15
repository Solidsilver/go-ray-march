set positional-arguments

start:
  go run ./cmd/image_render/image_render.go -t 10


start-prof:
  go run ./cmd/image_render/image_render.go -t 10 --cpuprofile=default.pgo

build-pgo:
  go run ./cmd/image_render/image_render.go -fov 10 -t 10 -d 3000 --cpuprofile=default.pgo
  go build -pgo=./default.pgo ./cmd/image_render/image_render.go

build-live-pgo:
  go run ./cmd/live_render/live_render.go -fov 10 -t 10 -d 3000 --cpuprofile=default.pgo
  go build -pgo=./default.pgo ./cmd/live_render/live_render.go

start-live:
  go run ./cmd/live_render/live_render.go -t 8 -d 1000