package chrome2

import (
	"context"
	"log"
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
		color.Red("Error occurred")
		if needFatal {
			log.Fatal(err)
		}
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
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
	err := chromedp.Run(ctxt, getStringsSlice(jsString, stringSlice))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetStringsSlice Error occured")
		color.Red("Error occured")
		if needFatal {
			log.Fatal(err)
		}
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
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
	err := chromedp.Run(ctxt, getString(jsString, resultString))
	if err != nil {
		utils.SendErrorToTelegram("CHROME: GetString Error occured")
		color.Red("Error occured")
		if needFatal {
			log.Fatal(err)
		}
	}
	if needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!")
	}
}
