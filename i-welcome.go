package main

import (
	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/messenger"
)

func GetStartedHandler(c *gbl.Context) {
	r := fb.CreateResponse(c)

	// Greeting text
	r.Image("https://upload.wikimedia.org/wikipedia/commons/d/da/Trivia_1.png")
	r.Text("Hello 👋 I am the trivia bot! My job is to ask you random trivia questions, would you like to do a round now?")

	// Quick reply buttons
	r.QR(fb.QRText("Yes!", "Yes"))
	r.QR(fb.QRText("No Thanks", "No"))

	// Add the start trivia prompt context so we know what to do when the user says yes or no
	bctx.Add(c, CStartTriviaPrompt, 1)
}
