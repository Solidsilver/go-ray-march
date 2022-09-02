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
	fovVert      float64
}

func NewCamera(pos utils.Vec3, sizeX int, sizeY int) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	// Camera{pos, utils.Vec3UnitX(), sizeX, sizeY, bgImg, 0}
	cam.Pos = pos
	cam.Dir = utils.Vec3UnitX()
	yP := utils.NewCopy(utils.Vec3UnitY())
	yP.Mult(0.12)
	cam.Dir.Add(*yP, cam.Dir)
	// cam.Dir.Add(utils.Vec3UnitY(), cam.Dir)
	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.frame = 0
	cam.Image = bgImg
	cam.centerOffset = Point{cam.SizeX / 2, cam.SizeY / 2}
	// cam.centerOffset = cam.centerOffset.Add(Point{int(cam.Pos.X)})
	// cam.w2 = cam.SizeX / 2
	// cam.h2 = cam.SizeY / 2
	cam.fov = 50
	cam.fovVert = (cam.fov / float64(sizeX)) * float64(sizeY)
	return cam
}

func (c *Camera) FlushToDisk() {
	imgName := fmt.Sprintf("/Users/solidsilver/Code/go-ray-march/rend_out/render%d.png", c.frame)
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
	ray.origin.Add(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	// return Ray{*rayPos, c.Dir}
}

func (c *Camera) RayForPixel2(px *Point, ray *Ray) {

	relPxPos := px.Sub(c.centerOffset)
	fovHalfRad := utils.DegToRad(c.fov / 2)
	adjX := float64(c.centerOffset.X) / math.Tan(fovHalfRad)
	vecX := utils.Vec3{X: float64(relPxPos.X), Y: adjX, Z: 0}
	vecX.Unit()

	fovYHalfRad := utils.DegToRad(c.fovVert / 2)
	adjY := float64(c.centerOffset.Y) / math.Tan(fovYHalfRad)
	vecY := utils.Vec3{X: float64(relPxPos.Y), Y: adjY, Z: 0}
	vecY.Unit()

	dirVec := utils.Vec3{X: 0, Y: vecX.X, Z: vecY.X}

	// rayPos := utils.NewAdd(utils.Vec3{X: 0, Y: float64(relPxPos.X), Z: float64(relPxPos.Y)}, c.Pos)
	ray.dir = *utils.NewAdd(c.Dir, dirVec)
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
