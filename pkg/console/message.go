package console

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/julienroland/usg"
)

func PrintSuccessMessage(msg string) {

	t := text.Colors{text.FgGreen}
	successMsg := fmt.Sprintf("%s %s\n", t.Sprintf("\n%s ", usg.Get.Tick), msg)
	fmt.Println(successMsg)
}

func PrintInfoMessage(msg string) {

	t := text.Colors{text.FgHiBlue}

	fmt.Println(t.Sprintf("%s ", usg.Get.Info), msg)
}

func PrintStatusMessage(msg string) {
	fmt.Println(fmt.Sprintf("\n%s", msg))
}

func PrintDefaultMessage(msg string) {
	fmt.Println(msg)
}

func PrintNextProcessMessage(msg string) {
	t := text.Colors{text.FgGreen}
	processMsg := fmt.Sprintf("%s %s", t.Sprintf("\n%s ", usg.Get.Play), msg)
	fmt.Println(processMsg)
}
