set positional-arguments

start:
  go run ./cmd/image_render/image_render.go -t 10

start-live:
  go run ./cmd/live_render/live_render.go -t 10