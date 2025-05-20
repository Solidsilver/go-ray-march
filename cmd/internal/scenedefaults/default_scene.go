package scenedefaults

import (
	"image/color"

	"github.com/Solidsilver/go-ray-march/pkg/drawables"
	"github.com/Solidsilver/go-ray-march/pkg/renderer"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

// NewDefaultRenderScene creates a new default scene with some objects and lights.
// This function is now located in the scenedefaults package.
func NewDefaultRenderScene(opts renderer.RenderOpts) *renderer.Renderer {
	// Setup Scene
	scene := renderer.NewBlankScene()
	scene.AddDrawables(
		// drawables.NewNamedSphere("s2", vec3.Vec3{X: 0, Y: 0, Z: 0}, 1.5, color.RGBA{255, 255, 255, 255}, false, false),
		drawables.NewMandelB(60, 1.5, 8, vec3.Zero, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewMandelB("m2", 60, 1.5, 12, vec3.Zero, color.RGBA{25, 35, 45, 255}, false),
		// drawables.NewCube(vec3.Vec3{X: 0, Y: 0, Z: 0}, 10, color.RGBA{237, 66, 22, 255}),
		// drawables.NewNamedTorus("t1", vec3.Vec3{X: 10, Y: -4, Z: -2}, 4, 0.25, color.RGBA{130, 156, 154, 255}),
		// drawables.NewNamedCube("b1", vec3.Vec3{X: -4, Y: -2, Z: -1.5}, 1, color.RGBA{255, 255, 255, 255}),
	)

	scene.AddLights(
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -15, Y: -1, Z: -1}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -2, Y: 2, Z: 0}, 1, color.RGBA{255, 255, 255, 255}, true, false),
		// drawables.NewLight(vec3.Vec3{X: 20, Y: -8, Z: -8}, 0.001, color.RGBA{100, 200, 200, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: 1, Y: 20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -20, Z: -8}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 3, Y: -7, Z: 8}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 101, Y: -20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 220, Y: -8, Z: -8}, 0.001, color.RGBA{100, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 10, Y: 20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -20, Z: -80}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 30, Y: -7, Z: 80}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1010, Y: -20, Z: 1}, 0.001, color.RGBA{199, 9, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 20, Y: -8, Z: -8}, 0.001, color.RGBA{4, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1, Y: -20, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -20, Z: -8}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 3, Y: -7, Z: 8}, 0.001, color.RGBA{200, 200, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 101, Y: -20, Z: 10}, 0.001, color.RGBA{199, 3, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -8, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 200, Y: -8, Z: -800}, 0.001, color.RGBA{100, 200, 1, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 10, Y: 120, Z: 10}, 0.001, color.RGBA{199, 219, 19, 255}, false),
		drawables.NewLight(vec3.Vec3{X: -80, Y: -20, Z: -80}, 0.001, color.RGBA{200, 19, 200, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 30, Y: -7, Z: 80}, 0.001, color.RGBA{200, 200, 9, 255}, false),
		drawables.NewLight(vec3.Vec3{X: 1010, Y: -20, Z: 1}, 0.001, color.RGBA{199, 55, 19, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: -180, Y: -100, Z: -8}, 0.001, color.RGBA{70, 80, 90, 255}, false),
		// drawables.NewNamedSphere("l2", vec3.Vec3{X: -15, Y: 8, Z: 8}, 1, color.RGBA{0, 255, 0, 255}, false),
		// drawables.NewLight(vec3.Vec3{X: -5, Y: -2, Z: 1}, 0.005, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l3", vec3.Vec3{X: -10, Y: -10, Z: 10}, 0.5, color.RGBA{69, 79, 79, 255}),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: -1, Z: -15}, 1, color.RGBA{240, 240, 240, 255}, false),
		// drawables.NewNamedSphere("l5", vec3.Vec3{X: -8, Y: -8, Z: -15}, 1, color.RGBA{200, 200, 200, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 1, Y: -1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: -1, Y: 1, Z: 15}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l1", vec3.Vec3{X: 0, Y: 0, Z: 17}, 1, color.RGBA{255, 255, 255, 255}, false),
		// drawables.NewNamedSphere("l5", vec3.Vec3{X: 8, Y: 8, Z: 15}, 1, color.RGBA{200, 200, 200, 255}, false),
	)

	cam := renderer.NewCameraFOV(vec3.Vec3{X: -10, Y: 0, Z: 0}, opts.DimX, opts.DimY, opts.Fov, opts.OutPath)

	cam.Up = vec3.UnitZ
	cam.Dir = vec3.UnitX

	// Create a new renderer instance using the constructor
	newRend := renderer.NewRenderer(scene, cam)

	return &newRend
}
