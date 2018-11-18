package main

import (
	"fmt"
	"html"
	"math/rand"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/messenger"
)

var qrMultipleChoice = []fb.QuickReply{
	fb.QRText("1st", "1st"),
	fb.QRText("2nd", "2nd"),
	fb.QRText("3rd", "3rd"),
	fb.QRText("4th", "4th"),
}

var qrTrueFalse = []fb.QuickReply{
	fb.QRText("True", "true"),
	fb.QRText("False", "false"),
}

func TriviaLoadQuestionsHandler(c *gbl.Context) {
	categories := c.GetStringSliceFlag(STriviaCategories)
	questionCount := c.GetIntFlag(SQuestionCount)

	// r := fb.CreateImmediateResponse(c)
	// r.RandomText(
	// 	"Hang tight, I'm thinking of some questions...",
	// 	"Hold on a sec, I'm just getting questions for you...",
	// )
	// r.Send()

	category := triviaAPI.GetCategoryFromString(categories[0])

	questions, err := triviaAPI.GetQuestions(questionCount, category)
	if err != nil {
		fmt.Println("Trivia Load ERR", err)
	}

	c.Flag(SQuestions, questions)
	c.Flag(SQuestionIDX, 0)
	c.Flag(SCorrectAnswers, make([]TriviaCorrectAnswer, 0))

	TriviaAskQuestionHandler(c)
}

func TriviaAskQuestionHandler(c *gbl.Context) {
	questions := c.GetFlag(SQuestions).([]TriviaQuestion)
	questionIndex := c.GetIntFlag(SQuestionIDX)

	r := fb.CreateResponse(c)

	if questionIndex == 0 {
		if len(questions) > 1 {
			r.Text("Let's get started! Here's your first question")
		} else {
			r.Text("Here's your question")
		}
	}

	if questionIndex >= len(questions) {
		FinishTriviaHandler(c)
		return
	}

	question := questions[questionIndex]

	if question.Type == "boolean" {
		TriviaAskTrueOrFalse(c, r, &question)
	} else {
		TriviaAskMultipleChoice(c, r, &question)
	}
}

func TriviaAskMultipleChoice(c *gbl.Context, r *fb.MBResponse, question *TriviaQuestion) {

	// r.TypingTime(1*time.Second, 0, 0, 0, 0)

	r.Text(html.UnescapeString(question.Question))

	var options = []string{
		question.CorrectAnswer,
		question.IncorrectAnswers[0],
		question.IncorrectAnswers[1],
		question.IncorrectAnswers[2],
	}

	correctAnswerIndex := 0

	rand.Shuffle(len(options), func(i, j int) {
		if j == correctAnswerIndex {
			correctAnswerIndex = i
		}

		options[i], options[j] = options[j], options[i]
	})

	for idx, option := range options {
		r.Text(fmt.Sprintf("%d. %s", idx+1, html.UnescapeString(option)))
	}

	r.QR(qrMultipleChoice...)

	correctAnswers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)
	correctAnswers = append(correctAnswers, TriviaCorrectAnswer{
		MultipleChoice: correctAnswerIndex,
	})
	c.Flag(SCorrectAnswers, correctAnswers)

	bctx.Add(c, CMPQuestionFollowup, 1)
}

func TriviaAskTrueOrFalse(c *gbl.Context, r *fb.MBResponse, question *TriviaQuestion) {
	r.Text("True or false?")
	r.Text(html.UnescapeString(question.Question))

	r.QR(qrTrueFalse...)

	var correct bool

	if question.CorrectAnswer == "True" {
		correct = true
	} else {
		correct = false
	}

	correctAnswers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)

	correctAnswers = append(correctAnswers, TriviaCorrectAnswer{
		TrueFalse: correct,
	})

	c.Flag(SCorrectAnswers, correctAnswers)

	bctx.Add(c, CTFQuestionFollowup, 1)
}
