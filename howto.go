package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"jaytaylor.com/html2text"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func get_question(words []string) ([]byte, error) {
	tag := words[0]
	question := strings.Join(words[1:], "+")
	url := fmt.Sprintf("https://api.stackexchange.com/2.2/search/advanced?q=%s&tagged=%s&accepted=true&order=desc&sort=votes&site=stackoverflow", question, tag)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytes, nil
}

func get_answer_id(bs []byte) (int, error) {
	var f map[string]interface{}
	err := json.Unmarshal(bs, &f)
	if err != nil {
		return 0, err
	}
	items := f["items"].([]interface{})
	if len(items) == 0 {
		return 0, errors.New("No answers!")
	}
	result := items[0].(map[string]interface{})
	id := result["accepted_answer_id"]
	return int(id.(float64)), nil
}

func get_answer(answer int) ([]byte, error) {
	url := fmt.Sprintf("https://api.stackexchange.com/2.2/answers/%d?order=desc&sort=activity&site=stackoverflow&filter=withbody", answer)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytes, nil
}

func get_answer_body(bs []byte) (string, error) {
	var f map[string]interface{}
	err := json.Unmarshal(bs, &f)
	if err != nil {
		return "", err
	}
	items := f["items"].([]interface{})
	if len(items) == 0 {
		return "", errors.New("No answers!")
	}
	result := items[0].(map[string]interface{})
	body := result["body"]
	return body.(string), nil
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal(errors.New("You have to pass a question!"))
	}
	bs, err := get_question(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	answer_id, err := get_answer_id(bs)
	if err != nil {
		log.Fatal(err)
	}
	bs, err = get_answer(answer_id)
	if err != nil {
		log.Fatal(err)
	}
	body, err := get_answer_body(bs)
	if err != nil {
		log.Fatal(err)
	}
	result, err := html2text.FromString(body, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
