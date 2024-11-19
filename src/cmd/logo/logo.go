package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	width  = 300
	height = 300
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 149, 182, 255}}, image.Point{}, draw.Src)

	for x := -width / 2; x < width/2; x++ {
		for y := -height / 2; y < height/2; y++ {
			if abs(x) < 100 {
				for z := -130; z <= -90; z++ {
					if abs(x)+z == y {
						img.Set(x+150, y+150, color.RGBA{255, 160, 0, 255})
					}
				}
				for z := -70; z <= -30; z++ {
					if abs(x)+z == y {
						img.Set(x+150, y+150, color.RGBA{255, 160, 0, 255})
					}
				}
				for z := -10; z <= 30; z++ {
					if abs(x)+z == y {
						img.Set(x+150, y+150, color.RGBA{255, 160, 0, 255})
					}
				}
				//for z := 70; z >= 30; z-- {
				//	if abs(x)+z == y {
				//		img.Set(x+150, y+150, color.Black)
				//	}
				//}
			}
			//if abs(x) < 200 && abs(y) < 200 {
			//	for z := -20; z < 20; z++ {
			//
			//	}
			//}
		}
	}

	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
