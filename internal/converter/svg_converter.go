package converter

import(
		"image"
		"image/png"
		"os"
		"fmt"

		"github.com/srwiley/oksvg"
		"github.com/srwiley/rasterx"
)

// converts SVG to image.Image
// image.Image is GO's universal img format that can be changed into wtv u want
func svgToImage(svgPath string) (image.Image, error){
	// Opens SVG file
	file, err := os.Open(svgPath)
	if err != nil {
		return nil, fmt.Errorf("Opening SVG file: %w", err)
	}
	defer file.Close()

	// Read file content
	icon, err := oksvg.ReadIconStream(file)
	if err != nil {
		return nil, fmt.Errorf("Reading SCG file: %w", err)
	}

	width := int(icon.ViewBox.W)
	height := int(icon.ViewBox.H)

	//defualt values if height and width unknown
	if width == 0 || height == 0 {
		width = 800
		height = 800
	}

	// Sets target rendering dimensions
	// Makes RGBA image "buffer"
	icon.SetTarget(0, 0, float64(width), float64(height))
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// Makes rasterizer and draws
	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	icon.Draw(dasher, 1)

	return rgba, nil
}

// converts SVG to PNG
func SVGToPNG(svgPath, outputPath string) error {
	// convert svg to img
	img, err := svgToImage(svgPath)
	if err != nil {
		return err
	}

	// create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Creating output file: %w", err)
	}
	defer outputFile.Close()

	//cencode as PNG 
	return png.Encode(outputFile, img)
}
