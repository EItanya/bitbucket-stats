package api

import (
	"fmt"

	"github.com/gosuri/uiprogress"
)

const prependStringLength = 30

func prependFormatFunc(f uiprogress.DecoratorFunc) uiprogress.DecoratorFunc {
	return func(b *uiprogress.Bar) string {
		prependedString := f(b)
		timeElapsed := b.TimeElapsedString()
		paddingSize := prependStringLength - (len(prependedString) + len(timeElapsed))
		padding := ""
		for i := 0; i < paddingSize; i++ {
			padding += " "
		}
		prependedString = fmt.Sprintf("%s:%s%s", prependedString, padding, timeElapsed)
		return prependedString
	}
}
