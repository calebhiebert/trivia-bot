package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/calebhiebert/gobbl"
	gobbldflow "github.com/calebhiebert/gobbl-dflow"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/messenger"
	"github.com/calebhiebert/gobbl/session"
	"github.com/joho/godotenv"
)

var triviaAPI *TriviaAPI

func main() {
	godotenv.Load(".env")

	keyfileB64 := os.Getenv("KEYFILE")

	if keyfileB64 == "" {
		panic("missing keyfile base64 env variable")
	}

	keyfileBytes, err := base64.StdEncoding.DecodeString(keyfileB64)
	if err != nil {
		panic(err)
	}

	triviaAPI = CreateTriviaAPI()
	rand.Seed(time.Now().UTC().UnixNano())

	dflowConfig, err := gobbldflow.LoadServiceKeyJSONBytes(keyfileBytes)
	dflow := gobbldflow.New(dflowConfig)

	gobblr := gbl.New()

	/*
		STANDARD MIDDLEWARE SETUP
		****************************************
	*/
	gobblr.Use(gbl.ResponderMiddleware())
	gobblr.Use(gbl.UserExtractionMiddleware())
	gobblr.Use(gbl.RequestExtractionMiddleware())
	gobblr.Use(fb.MarkSeenMiddleware())
	gobblr.Use(sess.Middleware(sess.MemoryStore()))
	gobblr.Use(bctx.Middleware())

	/*
		ROUTER SETUP
		****************************************
		Routers in this project are package-global and used here
	*/
	textRouter := gbl.TextRouter()
	ictxRouter := bctx.ContextIntentRouter()
	gobblr.Use(textRouter.Middleware())

	// LUIS is added at this point so that if any of our text routes match
	// we can skip the NLP process becuase we don't need to know the intent
	gobblr.Use(gobbldflow.Middleware(dflow))

	gobblr.Use(func(c *gbl.Context) {
		c.Infof("INTENT %v", c.GetFlag("intent"))
		fmt.Println(c.GetFlag("dflow"))
		c.Next()
	})

	gobblr.Use(ictxRouter.Middleware())

	// Fallback handler
	gobblr.Use(DefaultFallbackHandler)

	/*
		ROUTE SETUP
		****************************************
		All the project routes are defined here
	*/
	// Text Routes
	textRouter.Text("GET_STARTED", GetStartedHandler)
	textRouter.Text("PLAY_TRIVIA", TriviaBeginHandler)
	textRouter.Text("TRIVIA_DENY", TriviaStartDenyHandler)

	// Contextual Routes
	ictxRouter.All(bctx.I{"No"}, bctx.C{CStartTriviaPrompt}, TriviaStartDenyHandler)
	ictxRouter.All(bctx.I{"Yes"}, bctx.C{CStartTriviaPrompt}, TriviaBeginHandler)
	ictxRouter.NoContext(bctx.I{"Play"}, TriviaBeginHandler)

	// Category selector
	ictxRouter.All(bctx.I{"TriviaCategory"}, bctx.C{CStartTrivia}, TriviaCategorySelectedHandler)
	ictxRouter.FallbackAll(bctx.C{CStartTrivia}, TriviaCategorySelectionFallbackHandler)

	// Question count handler
	ictxRouter.All(bctx.I{"Number"}, bctx.C{CQuestionCount}, TriviaSelectQuestionCountHandler)
	ictxRouter.FallbackAll(bctx.C{CQuestionCount}, func(c *gbl.Context) { TriviaSelectQuestionCountFallbackHandler(c) })

	// Answer handlers
	ictxRouter.All(bctx.I{"Ordinal", "Number"}, bctx.C{CMPQuestionFollowup}, MultipleChoiceAnswerHandler)
	ictxRouter.FallbackAll(bctx.C{CMPQuestionFollowup}, MultipleChoiceAnswerHandlerFallback)

	ictxRouter.All(bctx.I{"True False"}, bctx.C{CTFQuestionFollowup}, TrueOrFalseAnswerHandler)
	ictxRouter.FallbackAll(bctx.C{CTFQuestionFollowup}, TrueOrFalseAnswerHandlerFallback)

	/*
		FACEBOOK MESSENGER SETUP
		****************************************
	*/
	mapi := fb.CreateMessengerAPI(os.Getenv("PAGE_ACCESS_TOKEN"))
	messengerIntegration := fb.MessengerIntegration{
		API:         mapi,
		VerifyToken: "frogs",
		DevMode:     true,
	}

	res, err := mapi.MessengerProfile(&fb.MessengerProfile{
		GetStarted: fb.GetStarted{
			Payload: "GET_STARTED",
		},
	})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Profile", res)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Server starting
	messengerIntegration.Listen(&http.Server{
		Addr: ":" + port,
	}, gobblr)
}
