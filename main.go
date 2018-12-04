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
	//"github.com/tidwall/gjson"
	//"bytes"

	"bytes"
	"github.com/tidwall/gjson"

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
func scratch_one_news(url string) {
	log.Printf("~~~~~~~~~~~~~~~~~~~~~scratch_one_news : %s" ,url)

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

	/*
		<div class="pg-rail-tall__body" itemprop="articleBody">
	*/


	htmlDoc.Find(".zn-body__paragraph").Each(func(i int, s *goquery.Selection){
		fmt.Println(s.Text())
	})
	txt := htmlDoc.Find(".zn-body__paragraph")
	article := ""
	log.Println(txt)
	for _,n := range txt.Nodes {
		if n.FirstChild != nil {
				article += n.FirstChild.Data
		}
	}
	log.Printf("item_uri-------------------->%s " ,article)


}

//func translate(article string ) (string, string){
//	google_web := "https://translate.google.cn/#view=home&op=translate&sl=en&tl=zh-CN&text=";
//	google_web = google_web + article;
//	res, err := http.Get(google_web)
//	if err != nil {
//		log.Println(err)
//		return "",""
//	}
//	defer res.Body.Close()
//	if res.StatusCode != 200 {
//		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
//		return "",""
//	}
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", "翻译出错"
//	}
//	// Load the HTML document
//
//
//	return "",""
//}


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
		//log.Println(destScript)


		size := bytes.Count([]byte(destScript),nil)
		log.Printf("%s page size=%d",gl_url,size)
		offset := strings.Index(destScript,"contentModel")
		log.Printf("contentModel offset=%d",offset)
		//log.Println(destScript[offset:(size-1)]);

		offset_1 := strings.Index(destScript[offset:size-1],"{")
		offset = offset+offset_1
		//log.Println(destScript[offset:size-1]);

		offset_1 = strings.Index(destScript[offset:size-1],"siblings")
		offset = offset+offset_1
		//log.Println(destScript[offset:size-1]);

		s := strings.Split(destScript[offset:size-1],"registryURL:")
		destScript = s[0];
		size = bytes.Count([]byte(destScript),nil)
		//log.Println(destScript)

		offset_1 = strings.Index(destScript,"{")
		offset = offset_1
		//log.Println(destScript[offset:size-1]);

		//now got you  "articleList"
		/*
			"articleList":{[
				{
				"uri": ...
				"headline": ...
				"thumbnail": ...
				"duration": ...
				"description": ...
				}
			]}

		*/

		result := gjson.Get(destScript[offset:size-1], "articleList")
		for _, item := range result.Array() {
			log.Println("----------------")
			//log.Println(item.String())
			item_uri :=item.Get("uri");
			//item_thumbnail :=item.Get("thumbnail");
			//item_headline :=item.Get("headline");
			//item_uri :=item.Get("uri");
			//item_description :=item.Get("description");
			log.Printf("item_uri--->%s " ,gl_url+item_uri.String())
			log.Printf("item_thumbnail--->%s " ,gl_url+item.Get("thumbnail").String())
			scratch_one_news(gl_url+item_uri.String())
			break;
		}


}


func main() {

	log.Print("build @ " + time.UnixDate)
	parse_cmdline();
	gl_url = "https://edition.cnn.com"

	init_logfile();
	log.Println("==============================Start:")
	scratch_website(gl_url)
	log.Println("==============================End:")
}
