package main

import (
	"fmt"
	"strconv"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/luis"
	"github.com/calebhiebert/gobbl/messenger"
)

var categoryQR = []fb.QuickReply{
	fb.QRText("General", "general"),
	fb.QRText("Books", "books"),
	fb.QRText("Movies", "movies"),
	fb.QRText("TV", "tv"),
	fb.QRText("Video Games", "video games"),
	fb.QRText("Science", "science"),
	fb.QRText("Mythology", "mythology"),
	fb.QRText("History", "history"),
	fb.QRText("Politics", "politics"),
}

var qnumQR = []fb.QuickReply{
	fb.QRText("1", "1"),
	fb.QRText("5", "5"),
	fb.QRText("10", "10"),
	fb.QRText("15", "15"),
}

func TriviaStartDenyHandler(c *gbl.Context) {
	r := fb.CreateResponse(c)

	r.Text("Sorry to hear that, just ask me to quiz you later if you change your mind!")
}

func TriviaBeginHandler(c *gbl.Context) {

	c.ClearFlag(SQuestionCount, SQuestionIDX, SQuestions, STriviaCategories)

	r := fb.CreateResponse(c)

	r.RandomText(
		"Sure, let's get this trivia started!",
		"👉😎👉 You got it, let's do this",
		"I've got you covered.",
	)
	r.RandomText(
		"Take a look at what I've got here, just pick one to get started",
		"What type of questions do you want?",
		"Step right up and choose a category!",
	)

	// Trivia category quick replies
	r.QR(categoryQR...)

	// Add correct context
	bctx.Add(c, CStartTrivia, 1)
}

func TriviaCategorySelectedHandler(c *gbl.Context) {
	r := fb.CreateResponse(c)

	luisResult := c.GetFlag("luis").(*luis.Response)

	categories := []string{}

	for _, entity := range luisResult.Entities {
		if entity.Resolution.Values != nil {
			for _, category := range entity.Resolution.Values {
				categories = append(categories, category)
			}
		}
	}

	// Redirect to the fallback if no categories were found
	if len(categories) == 0 {
		TriviaCategorySelectionFallbackHandler(c)
		return
	}

	r.RandomText(
		fmt.Sprintf("Good choice, I've always been partial to %s myself", categories[0]),
		fmt.Sprintf("%s eh? Let's do this", categories[0]),
		fmt.Sprintf("Questions about %s coming right up", categories[0]),
	)

	// Add the selected categoreies to the session
	c.Flag(STriviaCategories, categories)

	// Send the prompt for how many questions they would like
	r.RandomText(
		"How many questions?",
		"How much trivia do you want exactly?",
		"It would please me to know how many questions you desire",
	)

	// Add quick replies for number of questions
	r.QR(qnumQR...)

	bctx.Add(c, CQuestionCount, 1)
}

func TriviaCategorySelectionFallbackHandler(c *gbl.Context) {
	r := fb.CreateResponse(c)

	r.Text("Sorry, I didn't get that, I'm gonna need you to pick a category!")
	r.Text("You can pick one from the buttons below, or type one in")

	r.QR(categoryQR...)

	bctx.Add(c, CStartTrivia, 1)
}

func TriviaSelectQuestionCountHandler(c *gbl.Context) {
	luisResult := c.GetFlag("luis").(*luis.Response)

	if luisResult.Entities == nil || len(luisResult.Entities) == 0 {
		TriviaSelectQuestionCountFallbackHandler(c)
		return
	}

	var num = -1
	var err error

	// Extract the number result
	for _, entity := range luisResult.Entities {
		if entity.Type == "builtin.number" {
			num, err = strconv.Atoi(entity.Resolution.Value)
			if err == nil {
				break
			} else {
				fmt.Println("Number parse err", err)
			}
		}
	}

	// Make sure the number is actually a number
	if num < 0 {
		TriviaSelectQuestionCountFallbackHandler(c, "Unfortunately it's not possible to have a negative number of questions, sorry :(")
		return
	} else if num > 50 {
		TriviaSelectQuestionCountFallbackHandler(c, "Having more than 50 questions is not currently supported, sorry :(")
		return
	}

	c.Flag(SQuestionCount, num)

	TriviaLoadQuestionsHandler(c)
}

func TriviaSelectQuestionCountFallbackHandler(c *gbl.Context, alternateMessage ...string) {
	r := fb.CreateResponse(c)

	if alternateMessage != nil && len(alternateMessage) > 0 {
		for _, message := range alternateMessage {
			r.Text(message)
		}
	} else {
		r.Text("Sorry, I couldn't find a number in that")
	}

	r.Text("Please enter a number or choose from the options below")

	r.QR(qnumQR...)

	bctx.Add(c, CQuestionCount, 1)
}
