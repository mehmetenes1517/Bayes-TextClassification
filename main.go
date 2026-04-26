package main

import (
	bayes "classifier/Bayes"
	"fmt"
)

const (
	CRIME_FOLDER         = "./Datasets/Crime"
	ENTERTAINMENT_FOLDER = "./Datasets/Entertainment"
	POLITICS_FOLDER      = "./Datasets/Politics"
	SCIENCE_FOLDER       = "./Datasets/Science"
)

func EvaluateModel(crime map[string]float64, entertainment map[string]float64, politics map[string]float64, science map[string]float64) float64 {
	allcount := 0
	truecount := 0
	cr_files := bayes.GetFiles(CRIME_FOLDER)
	for _, el := range cr_files {
		if bayes.CheckFile(CRIME_FOLDER+"/"+el, "c", crime, entertainment, politics, science) {
			truecount++
		}
		allcount++
	}

	enter_files := bayes.GetFiles(ENTERTAINMENT_FOLDER)
	for _, el := range enter_files {
		if bayes.CheckFile(ENTERTAINMENT_FOLDER+"/"+el, "e", crime, entertainment, politics, science) {
			truecount++
		}
		allcount++
	}
	polit_files := bayes.GetFiles(POLITICS_FOLDER)
	for _, el := range polit_files {
		if bayes.CheckFile(POLITICS_FOLDER+"/"+el, "p", crime, entertainment, politics, science) {
			truecount++
		}
		allcount++
	}
	sci_files := bayes.GetFiles(SCIENCE_FOLDER)
	for _, el := range sci_files {
		if bayes.CheckFile(SCIENCE_FOLDER+"/"+el, "s", crime, entertainment, politics, science) {
			truecount++
		}
		allcount++
	}
	return float64(truecount) / float64(allcount)
}
func main() {

	crime_count := bayes.TokenizeWords(CRIME_FOLDER)
	entertainment_count := bayes.TokenizeWords(ENTERTAINMENT_FOLDER)
	politics_count := bayes.TokenizeWords(POLITICS_FOLDER)
	science_count := bayes.TokenizeWords(SCIENCE_FOLDER)

	merged_map := bayes.MergeMaps(bayes.MergeMaps(crime_count, entertainment_count), bayes.MergeMaps(politics_count, science_count))
	total_Wordcount := bayes.GetTotalWordSize(merged_map)

	P_Word := make(map[string]float64)
	for key, value := range merged_map {
		P_Word[key] = value / total_Wordcount
	} //P(Word)

	var total_EmailCount float64 = float64(len(bayes.GetContents(CRIME_FOLDER)) + len(bayes.GetContents(ENTERTAINMENT_FOLDER)) + len(bayes.GetContents(POLITICS_FOLDER)) + len(bayes.GetContents(SCIENCE_FOLDER)))
	P_Crime := float64(len(bayes.GetContents(CRIME_FOLDER))) / (total_EmailCount)
	P_Entertainment := float64(len(bayes.GetContents(ENTERTAINMENT_FOLDER))) / (total_EmailCount)
	P_Politics := float64(len(bayes.GetContents(POLITICS_FOLDER))) / (total_EmailCount)
	P_Science := float64(len(bayes.GetContents(SCIENCE_FOLDER))) / (total_EmailCount)
	//P(Class)

	P_Word_Crime := bayes.P_Word_Class(crime_count)
	P_Word_Entertainment := bayes.P_Word_Class(entertainment_count)
	P_Word_Politics := bayes.P_Word_Class(politics_count)
	P_Word_Science := bayes.P_Word_Class(science_count)
	//P(W | Class)

	//P(Class | Word)
	// P(Word | Class) * P(Class) / P/(Word)

	P_Crime_Word := bayes.P_Class_Word(P_Word_Crime, P_Crime, P_Word)
	P_Entertainment_Word := bayes.P_Class_Word(P_Word_Entertainment, P_Entertainment, P_Word)
	P_Politics_Word := bayes.P_Class_Word(P_Word_Politics, P_Politics, P_Word)
	P_Science_Word := bayes.P_Class_Word(P_Word_Science, P_Science, P_Word)

	fmt.Println("Accuracy : ", EvaluateModel(P_Crime_Word, P_Entertainment_Word, P_Politics_Word, P_Science_Word))
	/*

		 P(C | W)  = P(W | C) * P(C) / P(W)

		 P(C|W) -> probability of a class happening in all word set

		 P(W|C) -> probability of a word inside of the class set -----> crime_p[word]
		 																-> P_Word_Class(crime_tokens,"word")
		 P(C) -> total probability of a class				 -----> crim / crim + ent + poli + sci
		 																->P_Crime ,P_Ent,P_Politics,P_Science etc...
		 P(W) -> total probability of a word				 -----> crim[word] + ent[word] + poli[word] + sci[word] / crim + ent + poli + sci
		 																->bayes.P_Word(MergedWordMap,"word")

		Input = Email -> EmailTokens

		Argmax of these is classification result:
			Sum (Log(P(Crime | W)))
			Sum (Log(P(Entertainment | W)))
			Sum (Log(P(Politics | W)))
			Sum (Log(P(Science | W)))

	*/

}
