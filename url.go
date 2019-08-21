package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	var c chan string = make(chan string)

	//флаг
	wordPtr := flag.String("url", "foo", "a string")
	flag.Parse()


	//чтение из файла
	var urls string
	inputFile, err := os.Open(*wordPtr)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	defer inputFile.Close()
	data := make([]byte, 64)
	for{
		n, err := inputFile.Read(data)
		if err == io.EOF{
			break
		}
		urls += string(data[:n])
	}
	urlsArr := strings.Split(urls, " ")


	for _, url := range urlsArr[0:] {
		go parserChannel(c, url)
	}
	for _, url := range urlsArr[0:] {
		cha := <- c

		lastPos := strings.LastIndex(url, "//")
		fileName := url[lastPos+2:] + ".txt"

		file, err := os.Create(fileName)
		if err != nil{
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close()
		file.WriteString(cha)
	}

}


//конвертер байт кода в строку
func convert( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
		s[i] = string(b[i])
	}
	return strings.Join(s, "")
}


func parserChannel(c chan string, url string) {
	resp, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprint(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		c <- fmt.Sprintf("fetch: чтение %s: %v\n", url, err)
		return
	}
	c <- convert(b[:])
}
