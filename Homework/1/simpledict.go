package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DictRequestCaiyun struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictBasic struct {
	UsPhonetic string   `json:"us-phonetic"`
	Phonetic   string   `json:"phonetic"`
	UkPhonetic string   `json:"uk-phonetic"`
	UkSpeech   string   `json:"uk-speech"`
	UsSpeech   string   `json:"us-speech"`
	Explains   []string `json:"explains"`
}

type DictWeb struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type DictRespYoudao struct {
	ErrorCode    string                 `json:"errorCode"`
	Query        string                 `json:"query"`
	Translation  []string               `json:"translation"`
	Basic        DictBasic              `json:"basic"`
	Web          []DictWeb              `json:"web,omitempty"`
	Lang         string                 `json:"l"`
	Dict         map[string]interface{} `json:"dict,omitempty"`
	Webdict      map[string]interface{} `json:"webdict,omitempty"`
	TSpeakUrl    string                 `json:"tSpeakUrl,omitempty"`
	SpeakUrl     string                 `json:"speakUrl,omitempty"`
	ReturnPhrase []string               `json:"returnPhrase,omitempty"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func HexBuffToString(buff []byte) string {
	var ret string
	for _, value := range buff {
		str := strconv.FormatUint(uint64(value), 16)
		if len([]rune(str)) == 1 {
			ret = ret + "0" + str
		} else {
			ret = ret + str
		}
	}
	return ret
}

func queryyoudao(input string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	// 需申请appkey和secretkey
	appkey := "5ae912679d61e33a"
	secretkey := "drCDaiN64pOYYzIlu1ApI7sWkcl70C1c"

	u1 := uuid.NewV4()
	stamp := time.Now().Unix()
	instr := appkey + input + u1.String() + strconv.FormatInt(stamp, 10) + secretkey
	sig := sha256.Sum256([]byte(instr))
	var sigstr string = HexBuffToString(sig[:])
	data := make(url.Values, 0)
	data["q"] = []string{input}
	data["from"] = []string{"en"}
	data["to"] = []string{"zh-CSH"}
	data["appKey"] = []string{appkey}
	data["salt"] = []string{u1.String()}
	data["sign"] = []string{sigstr}
	data["signType"] = []string{"v3"}
	data["curtime"] = []string{strconv.FormatInt(stamp, 10)}

	resp, err := http.PostForm("https://openapi.youdao.com/api", data)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}

	var jsonObj DictRespYoudao
	json.Unmarshal(bodyText, &jsonObj)
	show(&jsonObj, os.Stdout)
}

func show(resp *DictRespYoudao, w io.Writer) {
	if resp.ErrorCode != "0" {
		fmt.Fprintln(w, "请输入正确的数据")
	}

	fmt.Fprintln(w, "[翻译]")
	for key, item := range resp.Translation {
		fmt.Fprintln(w, "\t", key+1, ".", item)
	}
	fmt.Fprintln(w, "[延伸]")
	for key, item := range resp.Basic.Explains {
		fmt.Fprintln(w, "\t", key+1, ".", item)
	}

	fmt.Fprintln(w, "[网络]")
	for key, item := range resp.Web {
		fmt.Fprintln(w, "\t", key+1, ".", item.Key)
		fmt.Fprint(w, "\t翻译:")
		for _, val := range item.Value {
			fmt.Fprint(w, val, ",")
		}
		fmt.Fprint(w, "\n")
	}
}

func querycaiyunai(word string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	client := &http.Client{}
	request := DictRequestCaiyun{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	fmt.Println(string(bodyText))
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println((dictResponse))
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	var wg sync.WaitGroup

	wg.Add(2) //有多少协程就添加多大的数
	go querycaiyunai(word, &wg)
	go queryyoudao(word, &wg)

	wg.Wait() // 等待
	fmt.Println("main exist")
}
