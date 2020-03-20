package utils

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/freetype"
)

func AddTextToImageByPth(srcPth string, drawX, drawY int, text, fntPth string, fntSize float64, rgba64 *color.RGBA64, dstPth string) error {

	//打开源文件
	fsrc, err := os.Open(srcPth)
	if err != nil {
		return err
	}

	//解码图片
	var imgSrc image.Image
	switch strings.ToLower(srcPth[len(srcPth)-3:]) {
	case "png":
		imgSrc, err = png.Decode(fsrc)
		if err != nil {
			return err
		}
	case "jpg":
		imgSrc, err = jpeg.Decode(fsrc)
		if err != nil {
			return err
		}
	default:
		return errors.New("不支持的文件格式")
	}

	//解析字体文件
	fd, err := ioutil.ReadFile(fntPth)
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(fd)
	if err != nil {
		return err
	}

	//添加文字
	canvas := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(canvas, canvas.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Src)
	c := freetype.NewContext()
	c.SetDst(canvas)
	c.SetClip(canvas.Bounds())
	c.SetSrc(image.NewUniform(rgba64))
	c.SetFont(f)
	c.SetFontSize(fntSize)
	_, err = c.DrawString(text, freetype.Pt(drawX, drawY))
	if err != nil {
		return err
	}

	//编码成png格式，输出到目标路径
	dst, err := os.Create(dstPth)
	if err != nil {
		return err
	}
	defer dst.Close()
	png.Encode(dst, canvas)

	return nil
}

func AddImageToImageByPth(bgPth string, t string, r io.Reader, drawX, drawY int, dstPth string) error {

	//编码成png格式，输出到目标路径
	var imgDg image.Image
	bg, err := os.Open(bgPth)
	if err != nil {
		return err
	}
	defer bg.Close()
	switch bgPth[len(bgPth)-3:] {
	case "png":
		imgDg, err = png.Decode(bg)
		if err != nil {
			return err
		}
	case "jpg":
		imgDg, err = jpeg.Decode(bg)
		if err != nil {
			return err
		}
	default:
		return errors.New("不支持的文件格式")
	}

	//解析源文件
	var imgSrc image.Image
	switch t {
	case "png":
		imgSrc, err = png.Decode(r)
		if err != nil {
			return err
		}
	case "jpg":
		imgSrc, err = jpeg.Decode(r)
		if err != nil {
			return err
		}
	default:
		return errors.New("不支持的文件格式")
	}

	//添加图片
	canvas := image.NewRGBA(imgDg.Bounds())
	draw.Draw(canvas, canvas.Bounds(), imgDg, imgDg.Bounds().Min, draw.Src)
	draw.Draw(canvas, imgSrc.Bounds().Add(image.Pt(drawX, drawY)), imgSrc, imgSrc.Bounds().Min, draw.Over)

	//编码成png格式，输出到目标路径
	dst, err := os.Create(dstPth)
	if err != nil {
		return err
	}
	defer dst.Close()
	png.Encode(dst, canvas)

	return nil
}

func AddTextToImageByReader(src io.Reader, fileType string, drawX, drawY int, text, fntPth string, fntSize float64, rgba64 *color.RGBA64) (*bytes.Buffer, error) {
	//解码图片
	var (
		imgSrc image.Image
		err    error
	)
	switch fileType {
	case "png":
		imgSrc, err = png.Decode(src)
		if err != nil {
			return nil, err
		}
	case "jpg":
		imgSrc, err = jpeg.Decode(src)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("不支持的文件格式")
	}

	//解析字体文件
	fd, err := ioutil.ReadFile(fntPth)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fd)
	if err != nil {
		return nil, err
	}

	//添加文字
	canvas := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(canvas, canvas.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Src)
	c := freetype.NewContext()
	c.SetDst(canvas)
	c.SetClip(canvas.Bounds())
	c.SetSrc(image.NewUniform(rgba64))
	c.SetFont(f)
	c.SetFontSize(fntSize)
	_, err = c.DrawString(text, freetype.Pt(drawX, drawY))
	if err != nil {
		return nil, err
	}

	//编码成png格式
	dst := bytes.Buffer{}
	png.Encode(&dst, canvas)

	return &dst, nil
}

func AddImageToImageByReader(bg io.Reader, t string, r io.Reader, drawX, drawY int) (*bytes.Buffer, error) {

	//编码成png格式，输出到目标路径
	var (
		imgDg image.Image
		err   error
	)
	switch t {
	case "png":
		imgDg, err = png.Decode(bg)
		if err != nil {
			return nil, err
		}
	case "jpg":
		imgDg, err = jpeg.Decode(bg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("不支持的文件格式")
	}

	//解析源文件
	var imgSrc image.Image
	switch t {
	case "png":
		imgSrc, err = png.Decode(r)
		if err != nil {
			return nil, err
		}
	case "jpg":
		imgSrc, err = jpeg.Decode(r)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("不支持的文件格式")
	}

	//添加图片
	canvas := image.NewRGBA(imgDg.Bounds())
	draw.Draw(canvas, canvas.Bounds(), imgDg, imgDg.Bounds().Min, draw.Src)
	draw.Draw(canvas, imgSrc.Bounds().Add(image.Pt(drawX, drawY)), imgSrc, imgSrc.Bounds().Min, draw.Over)

	//编码成png格式
	dst := bytes.Buffer{}
	png.Encode(&dst, canvas)

	return &dst, nil
}
