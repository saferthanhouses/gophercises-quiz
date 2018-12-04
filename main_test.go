package main

import "testing"


func TestAnswerComparisonWellFormed(t *testing.T){
	// it should accept well formed strings
	wfInput := "Buenos Aires"
	wfAnswer := "Buenos Aires"
	if got := compareAnswers(wfInput, wfAnswer); got != true {
		t.Errorf("compareAnswers(%s, %s) = %v, want %v", wfInput, wfAnswer, got, true)
	}
}

func TestAnswerComparisonBadlyFormed(t *testing.T){
	// it should accept well formed strings
	wfInput := "		Buenos Aires"
	wfAnswer := "Buenos Aires  "
	if got := compareAnswers(wfInput, wfAnswer); got != true {
		t.Errorf("compareAnswers(%s, %s) = %v, want %v", wfInput, wfAnswer, got, true)
	}
}

