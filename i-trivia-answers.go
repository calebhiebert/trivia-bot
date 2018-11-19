package main

import (
	"fmt"
	"html"
	"strconv"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/luis"
	"github.com/calebhiebert/gobbl/messenger"
)

var correctStrings = []string{
	"🎉🎉 Correct 🎉🎉", "✔️ You right", "Excellent 👌",
	"🎆🎇 YES 🎇🎆", "Oohhh, yeaaaahhhhhhhh 😎", "💯% correct",
	"💸💸💸 Yup 💸💸💸"}

var incorrectStrings = []string{"❌ Wrong", "😞 Incorrect", "Nope 😑"}

func MultipleChoiceAnswerHandler(c *gbl.Context) {
	luisResult := c.GetFlag("luis").(*luis.LUISResponse)

	if luisResult.Entities == nil || len(luisResult.Entities) == 0 {
		MultipleChoiceAnswerHandlerFallback(c)
		return
	}

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
		correctAnswers[len(correctAnswers)-1].UserAnsweredCorrectly = true
		r.RandomText(correctStrings...)
	} else {
		correctAnswers[len(correctAnswers)-1].UserAnsweredCorrectly = false
		r.RandomText(incorrectStrings...)
		r.Text("The correct answer is:\n" + html.UnescapeString(question.CorrectAnswer))
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

	choice := luisResult.Entities[0].Resolution.Values[0] == "true"

	r := fb.CreateImmediateResponse(c)

	correctAnswers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)
	correctAnswer := correctAnswers[len(correctAnswers)-1]

	if choice == correctAnswer.TrueFalse {
		correctAnswers[len(correctAnswers)-1].UserAnsweredCorrectly = true
		r.RandomText(correctStrings...)
	} else {
		correctAnswers[len(correctAnswers)-1].UserAnsweredCorrectly = false
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

	answers := c.GetFlag(SCorrectAnswers).([]TriviaCorrectAnswer)

	correctCount := 0

	for _, answer := range answers {
		if answer.UserAnsweredCorrectly {
			correctCount++
		}
	}

	r.RandomText(
		"That's it, let's see how you did",
		"All finished, calculating score...",
		"That's a wrap 🌯 Hang tight for your score",
	)

	r.Text("Drumroll please...")
	r.Text("🥁🥁🥁🥁🥁🥁🥁🥁🥁")

	percentCorrect := (float64(correctCount) / float64(len(answers))) * 100

	if percentCorrect == 0 {
		r.RandomText(
			"Yikes, you need to brush up on your trivia",
			"Maybe next time you'll get one right 😬",
		)
	} else if percentCorrect < 50 {
		r.RandomText(
			"Not the best score but at least you got some correct",
			"Oh dear, sadly not the worst I've seen",
		)
	} else if percentCorrect < 75 {
		r.RandomText(
			"Not bad!",
			"You're going places kid, good job",
		)
	} else if percentCorrect < 100 {
		r.RandomText(
			"Almost perfect! Maybe next time",
			"So close!",
		)
	} else if percentCorrect == 100 {
		r.RandomText(
			"🙌 A perfect score!",
			"Just fabulous 👌 100%!",
		)
	}

	r.Text(fmt.Sprintf("You got %d/%d", correctCount, len(answers)))

	r.QR(fb.QRText("Play Again", "PLAY_TRIVIA"))
}
