package main

import (
	"log"
	"time"
	"flag"
	"fmt"
	"os"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)
//"github.com/PuerkitoBio/goquery"
//"github.com/opesun/goquery"
var gl_url string
var gl_logfile string
func parse_cmdline(){


	flag.StringVar(&gl_url,"url","","web site")
	flag.StringVar(&gl_logfile,"logfile","","log file")
	flag.Parse();
	log.Print( flag.NFlag())
	//flag.Usage();

	fmt.Println("url=" + gl_url )
	fmt.Println("logfile=" + gl_logfile)

}
func init_logfile() {
	var (
		logFileName = flag.String("log", "log.txt", "Log file name")
	)
	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		log.Println("Fail to find", *logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
func scratch_website(url string){

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return;
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			return
		}
		// Load the HTML document
		htmlDoc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
			return
		}

		scripts := htmlDoc.Find("script")
		destScript := ""

		for _,n := range scripts.Nodes {
			if n.FirstChild != nil {
				if strings.Contains(n.FirstChild.Data,"articleList"){
					destScript = n.FirstChild.Data
					//log.Println(destScript)
					break;
				}
			}
		}

		offset := strings.Index(destScript,"contentModel")
		log.Println(destScript[offset])

}

func main() {
	log.Print("build @ " + time.UnixDate)
	parse_cmdline();
	init_logfile();
	log.Println("==============================Start:")
	scratch_website("https://edition.cnn.com")
	log.Println("==============================End:")
}
