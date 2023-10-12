set positional-arguments

start:
  # go run ./cmd/image_render/image_render.go -t 10 -fov 30 -d 30720x17280
  go run ./cmd/image_render/image_render.go -t 10 -fov 30 -d 2000

start-live:
  go run ./cmd/live_render/live_render.go -t 8 -fov 30 -d 2000
  # go run ./cmd/live_render/live_render.go -t 8 -fov 30 -d 3840x2160