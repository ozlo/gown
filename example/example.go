package main

import (
	"fmt"
	"github.com/ozlo/gown"
)

func main() {
	dictDir, _ := gown.GetWordNetDictDir()
	fmt.Printf("loading from %s\n", dictDir)
	wn, err := gown.LoadWordNet(dictDir)
	if err != nil {
		fmt.Printf("can't load WordNet from %v: %v\n", dictDir, err)
		return
	}

	printLookup(wn, "live")
	printLookupWithPartOfSpeech(wn, "computer", gown.POS_NOUN)
	printLookupWithPartOfSpeechAndSense(wn, "lemma", gown.POS_NOUN, 1)
	printLookupSensesWithPartOfSpeech(wn, "lemma", gown.POS_NOUN)

	words := []string { "computer", "computing machine", "computing device", "data processor", "electronic computer", "information processing system" }
	for _, word := range words {
		fmt.Printf("word: %s\n", word)
		printSenseIndexEntryAndSynset(wn, word, gown.POS_NOUN, 1)
	}

}

func printSenseIndexEntryAndSynset(wn *gown.WN, word string, pos int, senseId int) {
	senseIndexEntry := wn.LookupWithPartOfSpeechAndSense(word, pos, senseId)
	printSenseIndexEntry(wn, senseIndexEntry)
	printSynsetPtr(wn, senseIndexEntry.GetSynsetPtr())
}

func printLookupSensesWithPartOfSpeech(wn *gown.WN, word string, pos int) {
	fmt.Printf("\n===================\n\n")
	fmt.Printf("Lookup %q\n", word)
	for _, senseIndexEntry := range wn.LookupSensesWithPartOfSpeech(word, pos) {
		printSenseIndexEntry(wn, senseIndexEntry)
		fmt.Printf("\n")
	}
}

func printLookupWithPartOfSpeechAndSense(wn *gown.WN, word string, pos int, senseId int) {
	senseIndexEntry := wn.LookupWithPartOfSpeechAndSense(word, pos, senseId)
	printSenseIndexEntry(wn, senseIndexEntry)
}

func printLookup(wn *gown.WN, word string) {
	fmt.Printf("\n===================\n\n")
	fmt.Printf("Lookup %q\n", word)
	for _, senseIndexEntry := range wn.Lookup(word) {
		printSenseIndexEntry(wn, senseIndexEntry)
		fmt.Printf("\n")
	}
}

func printLookupWithPartOfSpeech(wn *gown.WN, word string, pos int) {
	fmt.Printf("\n===================\n\n")
	fmt.Printf("LookupWithPartOfSpeech %q %d\n", word, pos)
	dataIndexEntry := wn.LookupWithPartOfSpeech(word, pos)
	if dataIndexEntry == nil {
		fmt.Printf("Can't found a \"%s\" as a %s!\n", word, gown.PartOfSpeechToString(pos))
	} else {
		fmt.Printf("%s (%s)\n", word, gown.PartOfSpeechToString(pos))
		fmt.Printf("%v\n", *dataIndexEntry)
		for _, synsetOffset := range dataIndexEntry.GetSynsetOffsets() {
			printSynsetPtr(wn, wn.GetSynset(pos, synsetOffset))
			fmt.Printf("\n")
		}
	}
}

func printSenseIndexEntry(wn *gown.WN, senseIndexEntry *gown.SenseIndexEntry) {
	fmt.Printf("\t%s\n", senseIndexEntry.ToString())
}

func printSynsetPtr(wn *gown.WN, synsetPtr *gown.Synset) {
	if synsetPtr == nil {
		fmt.Printf("\tNO SYNSET!\n")
	} else {
		fmt.Printf("\tGloss: %s\n", synsetPtr.GetGloss())
		fmt.Printf("\tLexFile: %s POS: %s\n",
			synsetPtr.GetLexographerFilename(),
			gown.PartOfSpeechToString(synsetPtr.GetPartOfSpeech()))

		fmt.Printf("\twords:")
		for i, word := range synsetPtr.GetWords() {
			fmt.Printf(" %s (%d)", word, synsetPtr.GetLexIds()[i])
		}
		fmt.Printf("\n")

		fmt.Printf("\trelations:\n")
		for i, relation := range synsetPtr.GetRelationships() {
			printRelationship(wn, i, relation, synsetPtr.GetWords())
		}
	}
}

func printRelationship(wn *gown.WN, i int, relation gown.RelationshipEdge, srcWords []string) {
	fmt.Printf("\t\t%d: %s (%d) >> ", i, gown.RelationshipIdToString(relation.GetRelationshipType()), relation.GetRelationshipType())
	targetPtr := wn.GetSynset(relation.GetPartOfSpeech(), relation.GetSynsetOffset())
	if targetPtr != nil {
		srcWordNumber := relation.GetSourceWordNumber()
		if srcWordNumber > 0 {
			srcWordNumber-- // make it zero based
		}
		targetWordNumber := relation.GetTargetWordNumber()
		if targetWordNumber > 0 {
			targetWordNumber-- // make it zero based
		}

		star := ""
		if relation.GetSourceWordNumber() == 0 && relation.GetTargetWordNumber() == 0 {
			star = "*"
		}

		fmt.Printf("%s (%d) %v  (%d:%d:%s%s -> %d:%d:%s%s) {{%v}}\n",
			gown.PartOfSpeechToString(relation.GetPartOfSpeech()),
			relation.GetSynsetOffset(),
			targetPtr.GetWords(),

			relation.GetSourceWordNumber(),
			srcWordNumber,
			srcWords[srcWordNumber], star,

			relation.GetTargetWordNumber(),
			targetWordNumber,
			targetPtr.GetWords()[targetWordNumber], star,

			relation)
	} else {
		fmt.Printf("NIL RELATION\n")
	}
}
