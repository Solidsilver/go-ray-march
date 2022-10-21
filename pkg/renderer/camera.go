package renderer

import (
	"fmt"
	"image"
	"math"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Camera struct {
	Pos          utils.Vec3
	Dir          utils.Vec3
	SizeX        int
	SizeY        int
	Image        *image.RGBA
	frame        int
	centerOffset Point
	fov          float64
	fov_vRad     float64
	fov_hRad     float64
	aspect       float64
}

func NewCamera(pos utils.Vec3, sizeX int, sizeY int) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	cam.Pos = pos
	cam.Dir = utils.Vec3UnitX()
	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.aspect = float64(cam.SizeX) / float64(cam.SizeY)
	cam.frame = 0
	cam.Image = bgImg
	cam.centerOffset = Point{cam.SizeX / 2, cam.SizeY / 2}

	cam.fov = 10

	cam.fov_vRad = utils.DegToRad(cam.fov)
	cam.fov_hRad = math.Atan(math.Tan(cam.fov_vRad/2.0)*cam.aspect) * 2.0
	return cam
}

func NewCameraFOV(pos utils.Vec3, sizeX int, sizeY int, fov float64) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	cam.Pos = pos
	cam.Dir = utils.Vec3UnitX()
	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.aspect = float64(cam.SizeX) / float64(cam.SizeY)
	cam.frame = 0
	cam.Image = bgImg
	cam.centerOffset = Point{cam.SizeX / 2, cam.SizeY / 2}

	cam.fov = fov
	cam.fov_vRad = utils.DegToRad(cam.fov)
	cam.fov_hRad = math.Atan(math.Tan(cam.fov_vRad/2.0)*cam.aspect) * 2.0
	return cam
}

func (c *Camera) FlushToDisk() {
	imgName := fmt.Sprintf("../../rend_out/render%d.png", c.frame)
	c.frame = c.frame + 1
	utils.EncodePNGToPath(imgName, c.Image)

}

func (c Camera) NewRayForPixel(px Point) Ray {
	relPxPos := px.Sub(c.centerOffset)
	rayPos := utils.NewAdd(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	return Ray{*rayPos, c.Dir}
}

func (c *Camera) RayForPixel(px *Point, ray *Ray) {
	relPxPos := px.Sub(c.centerOffset)
	// rayPos := utils.NewAdd(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	ray.dir = c.Dir
	ray.origin = *utils.NewAdd(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	// return Ray{*rayPos, c.Dir}
}

func (c *Camera) RayForPixel2(px *Point, ray *Ray) {

	relPxPos := px.Sub(c.centerOffset)
	fovHalfRad := c.fov_hRad / 2
	adjX := float64(c.centerOffset.X) / math.Tan(fovHalfRad)
	vecX := utils.NewVec2(float64(relPxPos.X), adjX)
	vecX.Unit2()

	fovYHalfRad := c.fov_vRad / 2
	adjY := float64(c.centerOffset.Y) / math.Tan(fovYHalfRad)
	vecY := utils.NewVec2(float64(relPxPos.Y), adjY)
	vecY.Unit2()

	dirVec := utils.Vec3{X: 0, Y: vecX.X(), Z: vecY.X()}
	// dirVec.Cross(dirVec, c.Dir)

	// rayPos := utils.NewAdd(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	// ray.dir = *utils.NewAdd(c.Dir, dirVec)
	ray.dir.Add(c.Dir, dirVec)
	ray.origin = c.Pos
	// ray.origin.Add(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	// return Ray{*rayPos, c.Dir}
}

type Point struct {
	X int
	Y int
}

func (pt Point) Add(pt2 Point) Point {
	return Point{pt.X + pt2.X, pt.Y + pt2.Y}
}

func (pt Point) Sub(pt2 Point) Point {
	return Point{pt.X - pt2.X, pt.Y - pt2.Y}
}
