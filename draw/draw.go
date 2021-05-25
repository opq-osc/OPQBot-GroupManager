package draw

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math/big"

	"github.com/fogleman/gg"
)

//go:embed techno-hideo-1.ttf
var ttf []byte

func init() {
	ioutil.WriteFile("./techno-hideo-1.ttf", ttf, 0777)
}
func Draw6Number() ([]byte, string, error) {
	num := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		num += n.String()
	}

	c := gg.NewContext(300, 120)

	if err := c.LoadFontFace("./techno-hideo-1.ttf", 60); err != nil {
		panic(err)
	}
	c.MeasureString(num)
	c.SetHexColor("#FFFFFF")
	c.Clear()

	c.SetRGB(0, 0, 0)
	c.Fill()

	c.DrawStringWrapped(num, 20, 30, 0, 0, 300, 1.5, gg.AlignLeft)
	//c.SetRGB(0, 0, 0)
	//c.Fill()
	buf := new(bytes.Buffer)
	err := png.Encode(buf, c.Image())
	if err != nil {
		return nil, num, err
	}
	return buf.Bytes(), num, nil
}
