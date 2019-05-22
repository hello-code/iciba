package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Dict struct {
	XMLName     xml.Name `xml:"dict"`
	Text        string   `xml:",chardata"`
	Num         string   `xml:"num,attr"`
	ID          string   `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Key         string   `xml:"key"`
	Ps          []string `xml:"ps"`
	Pron        []string `xml:"pron"`
	Pos         []string `xml:"pos"`
	Acceptation []string `xml:"acceptation"`
	Sent        []struct {
		Text  string `xml:",chardata"`
		Orig  string `xml:"orig"`
		Trans string `xml:"trans"`
	} `xml:"sent"`
}

type Ps struct {
	XMLName xml.Name `xml:"ps"`
}

type Pron struct {
	XMLName xml.Name `xml:"pron"`
}

type Pos struct {
	XMLName xml.Name `xml:"pos"`
}

type Acceptation struct {
	XMLName xml.Name `xml:"acceptation"`
}

var counter int = 0
var APIKey string

func main() {
	result, err := os.Create("result.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()

	words, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer words.Close()

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		APIKey = "your key"
		url := "http://dict-co.iciba.com/api/dictionary.php?w=" + scanner.Text() + "&key=" + APIKey
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("GEt err:", err)
		}
		defer resp.Body.Close()

		// decoder
		dict := &Dict{}
		decoder := xml.NewDecoder(resp.Body)
		err = decoder.Decode(&dict)
		if err != nil {
			fmt.Println("decode err:", err)
		}

		var str strings.Builder
		str.WriteString(dict.Key + " ")

		if len(dict.Ps) == 0 {
			str.WriteString("☀[]☀[]")
		}
		if len(dict.Ps) == 1 {
			str.WriteString("☀[" + dict.Ps[0] + "]" + "☀")
		} else {
			for _, v := range dict.Ps { // phonetic symbol
				str.WriteString("☀[" + v + "] ")
			}
		}

		if len(dict.Pron) == 0 {
			str.WriteString("☀[]☀[]")
		}
		if len(dict.Pron) == 1 {
			str.WriteString("☀[" + dict.Pron[0] + "]" + "☀")
		} else {
			for _, v := range dict.Pron { // phonetic symbol
				str.WriteString("☀[" + v + "] ")
			}
		}

		if len(dict.Pos) == 0 {
			str.WriteString("☀☀☀☀☀☀☀☀☀☀☀☀")
		} else {
			var x int
			for i, v := range dict.Pos { // part of speech
				str.WriteString("☀" + v + "☀" + dict.Acceptation[i])
				x = i
			}
			if x < 6 {
				l := 6 - x
				j := 1
				for j < l {
					str.WriteString("☀" + "☀")
					x++
					j++
				}
			}
		}

		for i, v := range dict.Sent { // example sentence
			str.WriteString("☀" + v.Orig + "☀" + dict.Sent[i].Trans)
		}

		whole := strings.Replace(str.String(), "\n", "", -1)
		result.WriteString(whole + "\n")
		str.Reset()

		counter = counter + 1
		fmt.Printf("\r %d", counter)
		fmt.Printf("\033 %d", counter)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
