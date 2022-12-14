package renderer

import (
	"fmt"
	"image"
	"math"

	"github.com/Solidsilver/go-ray-march/pkg/utils"
	"github.com/Solidsilver/go-ray-march/pkg/vec3"
)

type Camera struct {
	Pos          vec3.Vec3
	Dir          vec3.Vec3
	SizeX        int
	SizeY        int
	Image        *image.RGBA
	frame        int
	centerOffset Point
	fov          float64
	fov_vRad     float64
	fov_hRad     float64
	aspect       float64
	up           vec3.Vec3
}

func NewCamera(pos vec3.Vec3, sizeX int, sizeY int) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	cam.Pos = pos
	cam.Dir = vec3.UnitX()
	cam.up = vec3.UnitZ()
	cam.SizeX = sizeX
	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.aspect = float64(cam.SizeX) / float64(cam.SizeY)
	cam.frame = 0
	cam.Image = bgImg
	cam.centerOffset = Point{cam.SizeX / 2, cam.SizeY / 2}

	cam.fov = 45

	cam.fov_vRad = utils.DegToRad(cam.fov)
	cam.fov_hRad = math.Atan(math.Tan(cam.fov_vRad/2.0)*cam.aspect) * 2.0
	return cam
}

func NewCameraFOV(pos vec3.Vec3, sizeX int, sizeY int, fov float64) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	cam.Pos = pos
	cam.Dir = vec3.UnitX()
	cam.up = vec3.UnitZ()
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
	imgName := fmt.Sprintf("../../rend_out/render%03d.png", c.frame)
	c.frame = c.frame + 1
	utils.EncodePNGToPath(imgName, c.Image)

}

func (c *Camera) RayForPixel(px Point) Ray {
	relPxPos := px.Sub(c.centerOffset)
	fovHalfRad := c.fov_hRad / 2
	adjX := float64(c.centerOffset.X) / math.Tan(fovHalfRad)
	vecX := vec3.Vec3{X: float64(relPxPos.X), Y: adjX, Z: 0}
	vecX = vecX.Unit()

	fovYHalfRad := c.fov_vRad / 2
	adjY := float64(c.centerOffset.Y) / math.Tan(fovYHalfRad)
	vecY := vec3.Vec3{X: float64(relPxPos.Y), Y: adjY, Z: 0}
	vecY = vecY.Unit()

	right := c.up.Cross(c.Dir).Mult(vecX.X)
	up := c.up.Mult(vecY.X)
	r2 := Ray{
		c.Pos,
		c.Dir.Add(up).Add(right),
	}

	return r2
}

func (c *Camera) Size() int64 {
	return int64(c.SizeX) * int64(c.SizeY)
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
