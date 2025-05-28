package ui

var colorPrefixes = map[string]string{
	"black":  "\033[30m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",

	"bg_black":  "\033[40m",
	"bg_red":    "\033[41m",
	"bg_green":  "\033[42m",
	"bg_yellow": "\033[43m",
	"bg_blue":   "\033[44m",
	"bg_purple": "\033[45m",
	"bg_cyan":   "\033[46m",
	"bg_white":  "\033[47m",

	"bright_black":  "\033[90m",
	"bright_red":    "\033[91m",
	"bright_green":  "\033[92m",
	"bright_yellow": "\033[93m",
	"bright_blue":   "\033[94m",
	"bright_purple": "\033[95m",
	"bright_cyan":   "\033[96m",
	"bright_white":  "\033[97m",

	"bg_bright_black":  "\033[100m",
	"bg_bright_red":    "\033[101m",
	"bg_bright_green":  "\033[102m",
	"bg_bright_yellow": "\033[103m",
	"bg_bright_blue":   "\033[104m",
	"bg_bright_purple": "\033[105m",
	"bg_bright_cyan":   "\033[106m",
	"bg_bright_white":  "\033[107m",

	"bold":          "\033[1m",
	"dim":           "\033[2m",
	"italic":        "\033[3m",
	"underline":     "\033[4m",
	"blink":         "\033[5m",
	"rapid_blink":   "\033[6m",
	"reverse":       "\033[7m",
	"hidden":        "\033[8m",
	"strikethrough": "\033[9m",
}

const reset = "\033[0m"

type colorString struct {
	text   string
	prefix string
}

func Color(str string) colorString {
	return colorString{text: str, prefix: ""}
}

func (c colorString) with(code string) colorString {
	return colorString{
		text:   c.text,
		prefix: c.prefix + code,
	}
}

// Color
func (c colorString) Black() colorString  { return c.with(colorPrefixes["black"]) }
func (c colorString) Red() colorString    { return c.with(colorPrefixes["red"]) }
func (c colorString) Green() colorString  { return c.with(colorPrefixes["green"]) }
func (c colorString) Yellow() colorString { return c.with(colorPrefixes["yellow"]) }
func (c colorString) Blue() colorString   { return c.with(colorPrefixes["blue"]) }
func (c colorString) Purple() colorString { return c.with(colorPrefixes["purple"]) }
func (c colorString) Cyan() colorString   { return c.with(colorPrefixes["cyan"]) }
func (c colorString) White() colorString  { return c.with(colorPrefixes["white"]) }

// Background color
func (c colorString) BackgroundBlack() colorString  { return c.with(colorPrefixes["bg_black"]) }
func (c colorString) BackgroundRed() colorString    { return c.with(colorPrefixes["bg_red"]) }
func (c colorString) BackgroundGreen() colorString  { return c.with(colorPrefixes["bg_green"]) }
func (c colorString) BackgroundYellow() colorString { return c.with(colorPrefixes["bg_yellow"]) }
func (c colorString) BackgroundBlue() colorString   { return c.with(colorPrefixes["bg_blue"]) }
func (c colorString) BackgroundPurple() colorString { return c.with(colorPrefixes["bg_purple"]) }
func (c colorString) BackgroundCyan() colorString   { return c.with(colorPrefixes["bg_cyan"]) }
func (c colorString) BackgroundWhite() colorString  { return c.with(colorPrefixes["bg_white"]) }

// Color bright
func (c colorString) BrightBlack() colorString  { return c.with(colorPrefixes["bright_black"]) }
func (c colorString) BrightRed() colorString    { return c.with(colorPrefixes["bright_red"]) }
func (c colorString) BrightGreen() colorString  { return c.with(colorPrefixes["bright_green"]) }
func (c colorString) BrightYellow() colorString { return c.with(colorPrefixes["bright_yellow"]) }
func (c colorString) BrightBlue() colorString   { return c.with(colorPrefixes["bright_blue"]) }
func (c colorString) BrightPurple() colorString { return c.with(colorPrefixes["bright_purple"]) }
func (c colorString) BrightCyan() colorString   { return c.with(colorPrefixes["bright_cyan"]) }
func (c colorString) BrightWhite() colorString  { return c.with(colorPrefixes["bright_white"]) }

// Background color bright
func (c colorString) BackgroundBrightBlack() colorString {
	return c.with(colorPrefixes["bg_bright_black"])
}
func (c colorString) BackgroundBrightRed() colorString {
	return c.with(colorPrefixes["bg_bright_red"])
}
func (c colorString) BackgroundBrightGreen() colorString {
	return c.with(colorPrefixes["bg_bright_green"])
}
func (c colorString) BackgroundBrightYellow() colorString {
	return c.with(colorPrefixes["bg_bright_yellow"])
}
func (c colorString) BackgroundBrightBlue() colorString {
	return c.with(colorPrefixes["bg_bright_blue"])
}
func (c colorString) BackgroundBrightPurple() colorString {
	return c.with(colorPrefixes["bg_bright_purple"])
}
func (c colorString) BackgroundBrightCyan() colorString {
	return c.with(colorPrefixes["bg_bright_cyan"])
}
func (c colorString) BackgroundBrightWhite() colorString {
	return c.with(colorPrefixes["bg_bright_white"])
}

// Style
func (c colorString) Bold() colorString          { return c.with(colorPrefixes["bold"]) }
func (c colorString) Dim() colorString           { return c.with(colorPrefixes["dim"]) }
func (c colorString) Italic() colorString        { return c.with(colorPrefixes["italic"]) }
func (c colorString) Underline() colorString     { return c.with(colorPrefixes["underline"]) }
func (c colorString) Blink() colorString         { return c.with(colorPrefixes["blink"]) }
func (c colorString) RapidBlink() colorString    { return c.with(colorPrefixes["rapid_blink"]) }
func (c colorString) Reverse() colorString       { return c.with(colorPrefixes["reverse"]) }
func (c colorString) Hidden() colorString        { return c.with(colorPrefixes["hidden"]) }
func (c colorString) Strikethrough() colorString { return c.with(colorPrefixes["strikethrough"]) }

// Output
func (c colorString) String() string {
	return c.prefix + c.text + reset
}
