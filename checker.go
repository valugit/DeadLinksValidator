package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var url string
var checked map[string]bool

type error interface {
	Error() string
}

func checkURL(link string) {

	if checked[link] == true {
		return
	}
	checked[link] = true

	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(errors.New("Can't open file"))
		log.Fatal(err)
	}
	defer file.Close()

	resp, err := http.Get(link)
	if err != nil {
		file.WriteString(link + " -> Error : " + err.Error() + "\n")
		file.Sync()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		file.WriteString(link + " -> OK\n")
		file.Sync()
	} else {
		file.WriteString(link + " -> Error " + strconv.Itoa(resp.StatusCode) + "\n")
		file.Sync()
	}

	check := regexp.MustCompile(url + `.*`)

	if check.MatchString(link) {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			file.WriteString(link + " -> Error : " + err.Error() + "\n")
			file.Sync()
			return
		}
		html := string(content)

		findURL(html)
	}
	return
}

func findURL(html string) {
	aurl := regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	rurl := regexp.MustCompile(`"\/([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])+"`)
	findaurl := aurl.FindAllString(html, -1)
	findrurl := rurl.FindAllString(html, -1)

	for _, value := range findaurl {
		checkURL(value)
	}

	for _, value := range findrurl {
		value = url + value[2:len(value)-1]
		checkURL(value)
	}
}

func main() {
	start := time.Now()
	printStart := []byte("Start : " + start.String() + "\n\n")
	ioutil.WriteFile("result.txt", printStart, 0644)
	reg := regexp.MustCompile(`https?://.*?/`)
	url = reg.FindString(os.Args[1])
	checked = make(map[string]bool)
	checkURL(url)

	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(errors.New("Can't open file"))
		log.Fatal(err)
	}
	defer file.Close()

	elapsed := time.Since(start)
	file.WriteString("\nUrls found : " + strconv.Itoa(len(checked)) + " | Time spent : " + elapsed.String())
	file.Sync()
}
