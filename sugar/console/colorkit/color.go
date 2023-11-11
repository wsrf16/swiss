package colorkit

import "strings"

var (
	GreenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	WhiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	YellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	RedBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	BlueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	MagentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	CyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Green        = string([]byte{27, 91, 51, 50, 109})
	White        = string([]byte{27, 91, 51, 55, 109})
	Yellow       = string([]byte{27, 91, 51, 51, 109})
	Red          = string([]byte{27, 91, 51, 49, 109})
	Blue         = string([]byte{27, 91, 51, 52, 109})
	Magenta      = string([]byte{27, 91, 51, 53, 109})
	Cyan         = string([]byte{27, 91, 51, 54, 109})
	Reset        = string([]byte{27, 91, 48, 109})
	DisableColor = false
)

func Spell(text string, colors ...string) string {
	pre := strings.Join(colors, "")
	return pre + text + Reset
}
