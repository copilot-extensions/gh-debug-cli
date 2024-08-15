package chat

import "fmt"

const (
	ColorDefault = "\x1b[37m"
	ColorRed     = "\x1b[31m"
	ColorGreen   = "\x1b[32m"
	ColorYellow  = "\x1b[33m"
	ColorMagenta = "\x1b[35m"
	ColorCyan    = "\x1b[36m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func yellow(s string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorDefault)
}

func magenta(s string) string {
	return fmt.Sprintf("%s%s%s", ColorMagenta, s, ColorDefault)
}

func cyan(s string) string {
	return fmt.Sprintf("%s%s%s", ColorCyan, s, ColorDefault)
}
