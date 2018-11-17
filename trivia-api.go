package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var CatGeneral TriviaCategory = 9
var CatBooks TriviaCategory = 10
var CatMovies TriviaCategory = 11
var CatTV TriviaCategory = 14
var CatVideoGames TriviaCategory = 15
var CatScience TriviaCategory = 17
var CatMythology TriviaCategory = 20
var CatHistory TriviaCategory = 23
var CatPolitics TriviaCategory = 24

type TriviaCategory int

type TriviaAPI struct {
	client *http.Client
}

type TriviaResponse struct {
	ResponseCode int              `json:"response_code"`
	Results      []TriviaQuestion `json:"results"`
}

type TriviaQuestion struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func CreateTriviaAPI() *TriviaAPI {
	return &TriviaAPI{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

/*
	GetQuestions will load a list of trivia questions from the API
*/
func (t *TriviaAPI) GetQuestions(count int, category TriviaCategory) ([]TriviaQuestion, error) {
	var jsonResponse TriviaResponse

	resp, err := t.client.Get(fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%d", count, category))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	if jsonResponse.ResponseCode != 0 {
		return nil, errors.New(fmt.Sprintf("API Returned incorrect response code %d", jsonResponse.ResponseCode))
	}

	return jsonResponse.Results, nil
}

func (t *TriviaAPI) GetCategoryFromString(category string) TriviaCategory {
	switch category {
	case "General":
		return CatGeneral
	case "Books":
		return CatBooks
	case "Movies":
		return CatMovies
	case "TV":
		return CatTV
	case "Video Games":
		return CatVideoGames
	case "Science":
		return CatScience
	case "Mythology":
		return CatMythology
	case "History":
		return CatHistory
	case "Politics":
		return CatPolitics
	}

	return 0
}
