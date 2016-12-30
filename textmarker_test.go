package resolveip

import (
	"fmt"
	"github.com/fatih/color"
	"testing"
)

func TestQuoting(t *testing.T) {
	fmt.Println(GenQuoter("»", "«")("Heiko"))
	fmt.Println(GenHighlighter(color.FgWhite, color.BgRed, color.Italic)("Heiko"))
}
