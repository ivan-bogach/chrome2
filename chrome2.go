package chrome2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ivan-bogach/nonsense"
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
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, openURL(url, &message)))
	if err != nil {
		statusCode, err := nonsense.SendStringToTelegram("CHROME: OpenURL Error occured: " + message)
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}

		color.Red("Error: %s", message)
		log.Fatal(err)
	}

	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
	}
}

func reload() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Reload(),
		chromedp.Sleep(5 * time.Second),
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: WaitVisible Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
		color.Red("Error in WaitVisible occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Reload and try again")
		err = chromedp.Run(ctxt, reload())
		if err != nil {
			log.Fatal(err)
		}
		WaitVisible(ctxt, selector, needLog, needFatal)
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
	}
}

func waitReady(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitReady(selector, chromedp.ByQuery),
	}
}

func WaitReady(ctxt context.Context, selector string, needLog, needFatal bool) {
	if needLog {
		c := color.New(color.FgGreen)
		c.Printf("Wait ready css:' %s ' - ", selector)
	}
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitReady(selector)))

	if err != nil {
		statusCode, err := nonsense.SendStringToTelegram("CHROME: WaitReady Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
		color.Red("Error in WaitReady occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Reload and try again")
		err = chromedp.Run(ctxt, reload())
		if err != nil {
			log.Fatal(err)
		}
		WaitReady(ctxt, selector, needLog, needFatal)
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: GetString Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: GetStringsSlice Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: GetReader Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: GetBool Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: Click Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
		color.Red("Error in Click occurred")
		if needFatal {
			log.Fatal(err)
		}
		color.Green("Try again")
		Click(ctxt, selector, needLog, needFatal)
	}
	err = chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, click(selector)))
	if err != nil {
		statusCode, err := nonsense.SendStringToTelegram("CHROME: Click Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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
		statusCode, err := nonsense.SendStringToTelegram("CHROME: SetInputValue Error occured")
		if err != nil {
			color.Red("can`t send error to telegram, error: ", err)
		}
		if statusCode != 200 {
			color.Red("can`t send error to telegram status code: ", statusCode)
		}
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

func WaitLoaded(ctxt context.Context) {
	var loaded bool
	GetBool(ctxt, `document.readyState !== 'ready' && document.readyState !== 'complete'`, &loaded, false, false)
	fmt.Print("Wait")
	n := 0
	for loaded {
		if n > 60 {
			utils.SendErrorToTelegram("Minute passed!")
		}
		fmt.Print(".")
		time.Sleep(1 * time.Second)
		GetBool(ctxt, `document.readyState !== 'ready' && document.readyState !== 'complete'`, &loaded, false, false)
		n++
	}
}

func parsePage(ctxt context.Context, js string) []string {
	var strSl []string
	GetStringsSlice(ctxt, js, &strSl, false, false)
	return strSl
}

func StringSliceFromPage(ctxt context.Context, url, js string, waitFor ...string) []string {
	OpenURL(ctxt, url, false)
	if len(waitFor) == 0 {
		WaitLoaded(ctxt)
	} else {
		for _, w := range waitFor {
			WaitVisible(ctxt, w, false, false)
		}
	}

	time.Sleep(1 * time.Second)
	newJSONSl := parsePage(ctxt, js)

	return newJSONSl
}
