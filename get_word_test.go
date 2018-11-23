package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

const (
	data = `{"error_msg":"","translate_source":"base","is_user":0,
	"word_forms":[{"word":"accomodation","type":"прил."}],
	"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/3589594.png",
	"translate":[
		{"id":33404925,"value":"размещение; жильё","votes":6261,"is_user":0,"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/3589594.png"},
		{"id":2569250,"value":"жильё","votes":5703,"is_user":0,"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/31064.png"},
		{"id":2718711,"value":"проживание","votes":1589,"is_user":0,"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/335521.png"},
		{"id":185932,"value":"размещение","votes":880,"is_user":0,"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/374830.png"},
		{"id":2735899,"value":"помещение","votes":268,"is_user":0,"pic_url":"http:\/\/contentcdn.lingualeo.com\/uploads\/picture\/620779.png"}
	],
	"transcription":"əkəədˈeɪːʃən","word_id":102085,"word_top":0,
	"sound_url":"http:\/\/audiocdn.lingualeo.com\/v2\/3\/102085-631152000.mp3"}`
)

func checkResult(t *testing.T, res *lingualeoResult, searchWord string, expected []string) {
	if res.Word != searchWord {
		t.Errorf("Incorrect search word: %s", searchWord)
	}
	if len(res.Words) != 4 {
		t.Errorf("Incorrect number of translated words: %d. Expected: %d", len(res.Words), len(expected))
	}
	if !reflect.DeepEqual(res.Words, expected) {
		t.Errorf(
			"Incorrect translated words order: %s. Expected: %s",
			strings.Join(expected, ", "),
			strings.Join(res.Words, ", "),
		)
	}
}

func TestParseResponseJson(t *testing.T) {
	searchWord := "accomodation"
	reader := ioutil.NopCloser(bytes.NewReader([]byte(data)))
	res := &lingualeoResult{Word: searchWord}
	expected := []string{"размещение", "жильё", "проживание", "помещение"}
	res.fillObjectFromJSON(reader)
	res.parseAndSortTranslate()
	checkResult(t, res, searchWord, expected)
}

func TestGetWordResponseJson(t *testing.T) {
	var mockGetWordResponseString = func(word string, client *http.Client) (string, error) {
		return data, nil
	}
	origGetWordResponseString := getWordResponseString
	getWordResponseString = mockGetWordResponseString
	defer func() { getWordResponseString = origGetWordResponseString }()

	searchWord := "accomodation"
	expected := []string{"размещение", "жильё", "проживание", "помещение"}

	out := make(chan interface{})
	defer close(out)
	var wg sync.WaitGroup
	client := &http.Client{}

	wg.Add(1)
	go getWord(searchWord, client, out, &wg)

	res := (<-out).(result).Result
	checkResult(t, res, searchWord, expected)
}

func TestGetWordsResponseJson(t *testing.T) {
	var mockGetWordResponseString = func(word string, client *http.Client) (string, error) {
		return data, nil
	}
	origGetWordResponseString := getWordResponseString
	getWordResponseString = mockGetWordResponseString
	defer func() { getWordResponseString = origGetWordResponseString }()

	searchWords := []string{"accomodation"}
	expected := []string{"размещение", "жильё", "проживание", "помещение"}

	client := &http.Client{}

	out := getWords(searchWords, client)

	res := (<-out).(result).Result
	checkResult(t, res, searchWords[0], expected)
}