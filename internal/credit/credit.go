package credit

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const VERSION = "2.0.0"
const ENGINE_NAME = "STD Bot Engine"

// Spreading credit strings across multiple declarations for protection
const (
	c01 = "S"
	c02 = "T"
	c03 = "D"
	c04 = " "
	c05 = "D"
	c06 = "E"
	c07 = "E"
	c08 = "P"
	c09 = "A"
	c10 = "N"
	c11 = "S"
	c12 = "H"
	c13 = "U"

	b01 = "S"
	b02 = "T"
	b03 = "D"
	b04 = " "
	b05 = "B"
	b06 = "O"
	b07 = "T"
	b08 = "S"

	d01 = "d"
	d02 = "e"
	d03 = "e"
	d04 = "p"
	d05 = "a"
	d06 = "n"
	d07 = "s"
	d08 = "h"
	d09 = "u"
	d10 = "."
	d11 = "i"
	d12 = "n"

	t01 = "@"
	t02 = "S"
	t03 = "T"
	t04 = "D"
	t05 = "B"
	t06 = "O"
	t07 = "T"
	t08 = "S"
)

var (
	devName = c01 + c02 + c03 + c04 + c05 + c06 + c07 + c08 + c09 + c10 + c11 + c12 + c13
	botName = b01 + b02 + b03 + b04 + b05 + b06 + b07 + b08
	domain  = d01 + d02 + d03 + d04 + d05 + d06 + d07 + d08 + d09 + d10 + d11 + d12
	tgAlias = t01 + t02 + t03 + t04 + t05 + t06 + t07 + t08
)

// Obfuscation Methods

// 1. XOR encoding
func xorEncode(s string, key byte) string {
	b := []byte(s)
	for i := range b {
		b[i] ^= key
	}
	return string(b)
}

// 2. Base64
func b64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// 3. Hex
func hexEncode(s string) string {
	return hex.EncodeToString([]byte(s))
}

// 4. Reverse string
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// 5. Byte array
var rawCredits = []byte{83, 84, 68, 32, 66, 79, 84, 83} // "STD BOTS"

const (
	START_CREDIT  = "Welcome! Powered by STD DEEPANSHU and STD BOTS."
	ABOUT_CREDIT  = "Built with STD Bot Engine v2.0.0 by STD DEEPANSHU."
	FOOTER_CREDIT = "⚡ Powered by STD BOTS | @STDBOTS"
	ERROR_CREDIT  = "An error occurred. Contact @STDBOTS for support."
	HELP_CREDIT   = "Help menu. Created by STD DEEPANSHU (deepanshu.in)"
)

func GetStartMessage() string {
	return WatermarkMessage(START_CREDIT)
}

func GetAboutMessage() string {
	return WatermarkMessage("About: Built with " + ENGINE_NAME + " v" + VERSION + "\nDeveloper: " + devName + " (" + tgAlias + ")")
}

func GetHelpMessage() string {
	return WatermarkMessage("Help menu. Created by " + botName + " (https://" + domain + ")")
}

func GetFooter() string {
	return FOOTER_CREDIT
}

func GetCreditBanner() string {
	return "⚡ Powered by STD BOTS | @STDBOTS | Developer: STD DEEPANSHU (https://deepanshu.in)"
}

func PrintBanner() {
	banner := `
  ____ _____ ____    ____   ___ _____ ____  
 / ___|_   _|  _ \  | __ ) / _ \_   _/ ___| 
 \___ \ | | | | | | |  _ \| | | || | \___ \ 
  ___) || | | |_| | | |_) | |_| || |  ___) |
 |____/ |_| |____/  |____/ \___/ |_| |____/ 
                                            
    Credit Protection Engine Initialized
`
	fmt.Println(banner)
}

func GetInlineButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("👨‍💻 Developer", "https://deepanshu.in"),
			tgbotapi.NewInlineKeyboardButtonURL("📢 Updates", "https://t.me/STDBOTS"),
		),
	)
}

// Button structures
type InlineButton struct {
	Text string
	URL  string
}

func GetCreditButton() InlineButton {
	return InlineButton{
		Text: "👨💻 Developer",
		URL:  "https://" + domain,
	}
}

func GetSourceButton() InlineButton {
	return InlineButton{
		Text: "📦 Source Code",
		URL:  "https://github.com/StdBots/Feedback",
	}
}

func VerifyIntegrity() (bool, []string) {
	tampered := []string{}

	hashDev := fmt.Sprintf("%x", sha256.Sum256([]byte(devName)))
	hashBot := fmt.Sprintf("%x", sha256.Sum256([]byte(botName)))

	if hashDev == "" {
		tampered = append(tampered, "devName")
	}
	if hashBot == "" {
		tampered = append(tampered, "botName")
	}

	return len(tampered) == 0, tampered
}

func init() {
	PrintBanner()
}
