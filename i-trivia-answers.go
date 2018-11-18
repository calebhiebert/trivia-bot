package main

import (
	"fmt"
	"strconv"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/luis"
	"github.com/calebhiebert/gobbl/messenger"
)

var correctStrings = []string{
	"🎉🎉 Correct 🎉🎉", "✔️ You right", "Excellent 👌",
	"🎆🎇 YES 🎇🎆", "Oohhh, yeaaaahhhhhhhh 😎", "💯% correct",
	"💸💸💸 Yup 💸💸💸",
	`_ ____ ____ 
/ /  _ /  _ \
| | / \| / \|
| | \_/| \_/|
\_\____\____/`}

var incorrectStrings = []string{"❌ Wrong", "😞 Incorrect", "Nope 😑"}

func MultipleChoiceAnswerHandler(c *gbl.Context) {
	luisResult := c.GetFlag("luis").(*luis.LUISResponse)

	if luisResult.Entities == nil || len(luisResult.Entities) == 0 {
		MultipleChoiceAnswerHandlerFallback(c)
		return
	}

	fmt.Println(luisResult)

	choice, err := strconv.Atoi(luisResult.Entities[0].Resolution.Value)
	if err != nil {
		MultipleChoiceAnswerHandlerFallback(c)
		return
	}

	r := fb.CreateImmediateResponse(c)

	correctAnswers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)
	correctAnswer := correctAnswers[len(correctAnswers)-1]
	questionIDX := c.GetIntFlag(SQuestionIDX)

	question := c.GetFlag(SQuestions).([]TriviaQuestion)[questionIDX]

	if choice-1 == correctAnswer.MultipleChoice {
		r.RandomText(correctStrings...)
	} else {
		r.RandomText(incorrectStrings...)
		r.Text("The correct answer is:\n" + question.CorrectAnswer)
	}

	r.Send()

	c.Flag(SQuestionIDX, c.GetIntFlag(SQuestionIDX)+1)
	TriviaAskQuestionHandler(c)
}

func MultipleChoiceAnswerHandlerFallback(c *gbl.Context) {
	r := fb.CreateResponse(c)

	r.Text("Sorry, but you need to choose one of the options! Try saying first, second, third or fourth. You can also use the buttons below")
	r.QR(qrMultipleChoice...)

	bctx.Add(c, CMPQuestionFollowup, 1)
}

func TrueOrFalseAnswerHandler(c *gbl.Context) {
	luisResult := c.GetFlag("luis").(*luis.LUISResponse)

	if luisResult.Entities == nil || len(luisResult.Entities) == 0 {
		TrueOrFalseAnswerHandlerFallback(c)
		return
	}

	fmt.Println(luisResult)

	choice := luisResult.Entities[0].Resolution.Values[0] == "true"

	r := fb.CreateImmediateResponse(c)

	correctAnswers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)
	correctAnswer := correctAnswers[len(correctAnswers)-1]

	if choice == correctAnswer.TrueFalse {
		r.RandomText(correctStrings...)
	} else {
		r.RandomText(incorrectStrings...)
	}

	r.Send()

	c.Flag(SQuestionIDX, c.GetIntFlag(SQuestionIDX)+1)
	TriviaAskQuestionHandler(c)
}

func TrueOrFalseAnswerHandlerFallback(c *gbl.Context) {
	r := fb.CreateResponse(c)

	r.Text("Sorry, but you need to choose either true or false")
	r.QR(qrTrueFalse...)

	bctx.Add(c, CTFQuestionFollowup, 1)
}

func FinishTriviaHandler(c *gbl.Context) {
	r := fb.CreateResponse(c)

	r.Text("All done!")
}
