package sprite

import "image"

// CreateRectangleGrid creates a two-dimensional slice of image.Rectangles
// that can be used to draw sprites from a sprite sheet.
//
// The rectangles are organized from left to right (columns), top to bottom (rows).
// Access them using cells[column][row] or cells[x][y], where:
// - 'x' corresponds to the column index (left to right)
// - 'y' corresponds to the row index (top to bottom).
func CreateRectangleGrid(cols, rows, width, height int) [][]image.Rectangle {
	cells := make([][]image.Rectangle, cols)

	for x := range cols {
		cells[x] = make([]image.Rectangle, rows)
		for y := range rows {
			cells[x][y] = image.Rect(
				x*width,
				y*height,
				(x*width)+width,
				(y*height)+height,
			)
		}
	}

	return cells
}
