package resolveip

import (
    "testing"
    "fmt"
    "github.com/fatih/color"
)

func TestQuoting(t *testing.T) {
    fmt.Println(GenQuoter("»", "«")("Heiko"))
    fmt.Println(GenHighlighter(color.FgWhite, color.BgRed, color.Italic)("Heiko"))
}
