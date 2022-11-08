package main

import (
	_ "image/jpeg"
	"time"

	"github.com/cnjack/epaper"
)

func main() {
	epaper.Init(400, 300)
	defer epaper.Close()
	epaper.Clear()
	// f, err := os.Open("./demo.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// img, _, err := image.Decode(f)
	// if err != nil {
	// 	panic(err)
	// }
	// epaper.DrawData(epaper.GetBuffer(img))
	time.Sleep(5 * time.Second)

}
