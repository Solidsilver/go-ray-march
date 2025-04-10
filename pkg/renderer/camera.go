package renderer

import (
	"fmt"
	"image"
	"math"
	"os"
	"sync"

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
	flushDir     string
	flushMtx     sync.Mutex
}

type CameraOpts struct {
	Position vec3.Vec3
	Size     utils.Vec2[float64]
	Fov      float64
	ImgDir   string
}

func NewCamera(pos vec3.Vec3, sizeX int, sizeY int, imgOut string) *Camera {
	return NewCameraFOV(pos, sizeX, sizeY, 45, imgOut)
}

func NewCameraFOV(pos vec3.Vec3, sizeX int, sizeY int, fov float64, imgOut string) *Camera {
	opts := CameraOpts{
		pos,
		utils.NewVec2(float64(sizeX), float64(sizeY)),
		fov,
		imgOut,
	}
	return NewCameraOpts(opts)

}

func NewCameraOpts(opts CameraOpts) *Camera {
	sizeX := int(opts.Size.X())
	sizeY := int(opts.Size.Y())
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	// bgImg := image.NewRGBA(image.Rect(0, 0, 0, 0))
	cam := new(Camera)
	cam.Pos = opts.Position
	cam.Dir = vec3.UnitX
	cam.up = vec3.UnitZ

	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.aspect = float64(cam.SizeX) / float64(cam.SizeY)
	cam.frame = 0
	cam.Image = bgImg
	cam.centerOffset = Point{cam.SizeX / 2, cam.SizeY / 2}

	cam.fov = opts.Fov
	cam.fov_vRad = utils.DegToRad(cam.fov)
	cam.fov_hRad = math.Atan(math.Tan(cam.fov_vRad/2.0)*cam.aspect) * 2.0

	cam.flushDir = opts.ImgDir
	return cam
}

// Dim returns the dimensions of the camera in pixels
func (c *Camera) Dim() utils.Vec2[int] {
	return utils.NewVec2(c.SizeX, c.SizeY)
}

func (c *Camera) FlushToDisk() {
	os.Mkdir(c.flushDir, os.ModePerm)
	c.flushMtx.Lock()
	curFrame := c.frame
	c.frame += 1
	c.flushMtx.Unlock()
	imgName := fmt.Sprintf("%s/render%03d.png", c.flushDir, curFrame)
	// log.Info().Msgf("Encoding to path: %s", imgName)
	utils.EncodePNGToPath(imgName, c.Image)
	// log.Info().Msgf("Encoded to path: %s", imgName)
}

func (c *Camera) Reset() {
	c.Image = image.NewRGBA(image.Rect(0, 0, c.SizeX, c.SizeY))
	// make the entire image black
	// for x := 0; x < c.SizeX; x++ {
	// 	for y := 0; y < c.SizeY; y++ {
	// 		c.Image.Set(x, y, color.RGBA{0, 0, 0, 0})
	// 	}
	// }
}

func (c *Camera) GetBytes() ([]byte, error) {
	return utils.EncodeImageToBytes(c.Image, utils.IMG_PNG)
}

func (c *Camera) RayForPixel(px Point) Ray {
	relPxPos := px.Sub(c.centerOffset)
	fovHalfRad := c.fov_hRad / 2
	adjX := float64(c.centerOffset.X) / math.Tan(fovHalfRad)
	vecX := vec3.Vec3{X: float64(relPxPos.X), Y: adjX, Z: 0}
	vecX = vecX.ToUnit()

	fovYHalfRad := c.fov_vRad / 2
	adjY := float64(c.centerOffset.Y) / math.Tan(fovYHalfRad)
	vecY := vec3.Vec3{X: float64(relPxPos.Y), Y: adjY, Z: 0}
	vecY = vecY.ToUnit()

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

func (c *Camera) MoveUp(amt float64) {
	c.Pos = c.Pos.Sub(c.up.Mult(amt))
}

func (c *Camera) MoveDown(amt float64) {
	c.Pos = c.Pos.Add(c.up.Mult(amt))
}

func (c *Camera) MoveLeft(amt float64) {
	c.Pos = c.Pos.Add(c.Dir.Cross(c.up).Mult(amt))
}

func (c *Camera) MoveRight(amt float64) {
	c.Pos = c.Pos.Sub(c.Dir.Cross(c.up).Mult(amt))
}

func (c *Camera) MoveForward(amt float64) {
	c.Pos = c.Pos.Add(c.Dir.Mult(amt))
}

func (c *Camera) MoveBackward(amt float64) {
	c.Pos = c.Pos.Sub(c.Dir.Mult(amt))
}

func (c *Camera) RotateLeft(amt float64) {
	c.Dir = c.Dir.Mult(1 - amt).Add(c.up.Cross(c.Dir).Mult(amt))
}

func (c *Camera) RotateRight(amt float64) {
	c.Dir = c.Dir.Mult(1 - amt).Sub(c.up.Cross(c.Dir).Mult(amt))
}
