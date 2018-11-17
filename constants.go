package main

var CStartTriviaPrompt string = "StartTriviaPrompt"
var CStartTrivia string = "StartTrivia"
var CQuestionCount string = "QuestionCount"
var CTFQuestionFollowup string = "TFQuestionFollowup"
var CMPQuestionFollowup string = "MPQuestionFollowup"

/*
	SESSION VARIABLE NAMES
	It can be useful to keep track of what your session variables are
	and what types they contain
*/

// Which categories the user has chosen
// type []string
var STriviaCategories = "sess:trivia-categories"

// The number of questions the user would like
// type int
var SQuestionCount = "sess:question-count"

// The trivia questions being asked to this user
// type []TriviaQuestion
var SQuestions = "sess:trivia-questions"

// The index of the question the user is currently on
// type int
var SQuestionIDX = "sess:trivia-question-idx"

// A slice storing the correct answers for each question
// type []TriviaCorrectAnswer
var SCorrectAnswers = "sess:trivia-correct-answers"

type TriviaCorrectAnswer struct {
	TrueFalse      bool
	MultipleChoice int
}
