package renderer

import (
	"fmt"
	"image"

	"solidsilver.dev/go-ray-marching/pkg/utils"
)

type Camera struct {
	Pos   utils.Vec3
	Dir   utils.Vec3
	SizeX int
	SizeY int
	Image *image.RGBA
	frame int
}

func NewCamera(pos utils.Vec3, sizeX int, sizeY int) *Camera {
	bgImg := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	cam := new(Camera)
	// Camera{pos, utils.Vec3UnitX(), sizeX, sizeY, bgImg, 0}
	cam.Pos = pos
	cam.Dir = utils.Vec3UnitX()
	// cam.Dir.Add(utils.Vec3UnitY(), cam.Dir)
	cam.SizeX = sizeX
	cam.SizeY = sizeY
	cam.frame = 0
	cam.Image = bgImg
	return cam
}

func (c *Camera) FlushToDisk() {
	imgName := fmt.Sprintf("/Users/solidsilver/Code/go-ray-march/rend_out/render%d.png", c.frame)
	c.frame = c.frame + 1
	utils.EncodePNGToPath(imgName, c.Image)

}

type Point struct {
	X int
	Y int
}
