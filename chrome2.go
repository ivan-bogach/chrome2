package chrome2

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ivan-bogach/utils"
	"github.com/knq/chromedp"
)

func InitHeadLess(pathForUserDataDir string) (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{chromedp.Flag("no-sandbox", true), chromedp.Flag("headless", true), chromedp.Flag("disable-gpu", true)}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(allocContext)
}

func Init(pathForUserDataDir string) (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{chromedp.Flag("no-sandbox", true)}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(allocContext)
}
func InitWithProxy(pathForUserDataDir, proxyName string) (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{chromedp.Flag("no-sandbox", true), chromedp.ProxyServer(proxyName)}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(allocContext)
}
func openURL(url string, message *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptOpenURL(url), message),
	}
}

func OpenURL(ctxt context.Context, url string, needLog bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Opening page url %s - ", url)
	}

	var message string
	err := chromedp.Run(ctxt, openURL(url, &message))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: OpenURL Error occured: " + message)
		color.Red("Error: %s", message)
		log.Fatal(err)
	}

	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
	}
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func waitVisible(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	}
}

func WaitVisible(ctxt context.Context, selector string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Wait visible css:' %s ' - ", selector)
	}
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitVisible(selector)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: WaitVisible Error occured")
		color.Red("Error in WaitVisible occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		WaitVisible(ctxt, selector, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
	}
}

func getString(jsString string, resultString *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetString(jsString), resultString),
	}
}

func GetString(ctxt context.Context, jsString string, resultString *string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Getting a string ' %s  ' - ", jsString)
	}
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getString(jsString, resultString)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetString Error occured")
		color.Red("Error in GetString occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		GetString(ctxt, jsString, resultString, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}

func getStringsSlice(jsString string, resultSlice *[]string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetStringsSlice(jsString), resultSlice),
	}
}

func GetStringsSlice(ctxt context.Context, jsString string, stringSlice *[]string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Getting a strings slice ' %s  ' - ", jsString)
	}
	color.Green("")
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getStringsSlice(jsString, stringSlice)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetStringsSlice Error occured")
		color.Red("Error in GetStringsSlice occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		GetStringsSlice(ctxt, jsString, stringSlice, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}

func GetReader(ctxt context.Context, jsString string, needLog, needFatal bool) *strings.Reader {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Getting a string ' %s  ' - ", jsString)
	}
	var resultString string
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getString(jsString, &resultString)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetReader Error occured")
		color.Red("Error in GetReader occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		GetReader(ctxt, jsString, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
	return strings.NewReader(resultString)
}

func getBool(jsBool string, resultBool *bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetBool(jsBool), resultBool),
	}
}

func GetBool(ctxt context.Context, jsBool string, resultBool *bool, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Getting a string ' %s  ' - ", jsBool)
	}
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getBool(jsBool, resultBool)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetBool Error occured")
		color.Red("Error in GetBool occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		GetBool(ctxt, jsBool, resultBool, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}

func click(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(selector, chromedp.ByQuery),
	}
}

func Click(ctxt context.Context, selector string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Click selector: ' %s '  - ", selector)
	}
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitVisible(selector)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: Click Error occured")
		color.Red("Error in Click occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		Click(ctxt, selector, needLog, needFatal)
	}
	err = chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, click(selector)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: Click Error occured")
		color.Red("Error in Click occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		Click(ctxt, selector, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}

func setInputValue(selector, value string, result *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptSetInputValue(selector, value), result),
	}
}

func SetInputValue(ctxt context.Context, selector, value string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Setting an input >>>%s<<< value - >>>%s<<<", selector, value)
	}
	color.Green("")
	var resultOperation string

	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, setInputValue(selector, value, &resultOperation)))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: SetInputValue Error occured")
		color.Red("Error in SetInputValue occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		SetInputValue(ctxt, selector, value, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}
