package main

import (
	"fmt"
	"github.com/ozlo/gown"
)

func main() {
	dictDir, _ := gown.GetWordNetDictDir()
	wn, err := gown.LoadWordNet(dictDir)
	if err != nil {
		fmt.Printf("can't load WordNet from %v: %v\n", dictDir, err)
		return
	}

	printLookup(wn, "live")
	printLookupWithPartOfSpeech(wn, "computer", gown.POS_NOUN)
}

func printLookup(wn *gown.WN, word string) {
	fmt.Printf("\n===================\n\n")
	fmt.Printf("live on\n")
	for resultId, senseIndexEntry := range wn.Lookup(word) {
		printSenseIndexEntry(wn, resultId, senseIndexEntry)
		fmt.Printf("\n")
	}
}

func printLookupWithPartOfSpeech(wn *gown.WN, word string, pos int) {
	fmt.Printf("\n===================\n\n")
	fmt.Printf("live on\n")
	dataIndexEntry := wn.LookupWithPartOfSpeech(word, pos)
	if dataIndexEntry == nil {
		fmt.Printf("Can't found a \"%s\" as a %s!\n", word, gown.PART_OF_SPEECH_ID_TO_STRING[pos])
	} else {
		fmt.Printf("%s (%s)\n", word, gown.PART_OF_SPEECH_ID_TO_STRING[pos])
		fmt.Printf("%v\n", *dataIndexEntry)
		for _, synsetOffset := range dataIndexEntry.SynsetOffsets {
			printSynsetPtr(wn, wn.GetSynset(pos, synsetOffset))
			fmt.Printf("\n")
		}
	}
}

func printSenseIndexEntry(wn *gown.WN, resultId int, senseIndexEntry *gown.SenseIndexEntry) {
	fmt.Printf("\tresultId: %2d POS: (%d) %s TagCount: %d synsetOffset: %d\n",
		resultId,
		senseIndexEntry.PartOfSpeech,
		gown.PART_OF_SPEECH_ID_TO_STRING[senseIndexEntry.PartOfSpeech],
		senseIndexEntry.TagCount,
		senseIndexEntry.SynsetOffset,
	)
	synsetPtr := wn.GetSynset(senseIndexEntry.PartOfSpeech, senseIndexEntry.SynsetOffset)
	printSynsetPtr(wn, synsetPtr)
}

func printSynsetPtr(wn *gown.WN, synsetPtr *gown.Synset) {
	if synsetPtr == nil {
		fmt.Printf("\tNO SYNSET!\n")
	} else {
		fmt.Printf("\tGloss: %s\n", synsetPtr.Gloss)
		fmt.Printf("\tLexFile: %s POS: %s\n",
			gown.LEXOGRAPHER_FILE_NUM_TO_NAME[synsetPtr.LexographerFilenum],
			gown.PART_OF_SPEECH_ID_TO_STRING[synsetPtr.PartOfSpeech])

		fmt.Printf("\twords:")
		for i, word := range synsetPtr.Words {
			fmt.Printf(" %s (%d)", word, synsetPtr.LexIds[i])
		}
		fmt.Printf("\n")

		fmt.Printf("\trelations:\n")
		for i, relation := range synsetPtr.Relationships {
			printRelationship(wn, i, relation, synsetPtr.Words)
		}
	}
}

func printRelationship(wn *gown.WN, i int, relation gown.RelationshipEdge, srcWords []string) {
	fmt.Printf("\t\t%d: %s (%d) >> ", i, gown.RELATIONSHIP_ID_TO_STRING[relation.RelationshipType], relation.RelationshipType)
	targetPtr := wn.GetSynset(relation.PartOfSpeech, relation.SynsetOffset)
	if targetPtr != nil {
		srcWordNumber := relation.SourceWordNumber
		if srcWordNumber > 0 {
			srcWordNumber-- // make it zero based
		}
		targetWordNumber := relation.TargetWordNumber
		if targetWordNumber > 0 {
			targetWordNumber-- // make it zero based
		}

		star := ""
		if relation.SourceWordNumber == 0 && relation.TargetWordNumber == 0 {
			star = "*"
		}

		fmt.Printf("%s (%d) %v  (%d:%d:%s%s -> %d:%d:%s%s) {{%v}}\n",
			gown.PART_OF_SPEECH_ID_TO_STRING[relation.PartOfSpeech],
			relation.SynsetOffset,
			targetPtr.Words,
			relation.SourceWordNumber,
			srcWordNumber,
			srcWords[srcWordNumber], star,
			relation.TargetWordNumber,
			targetWordNumber,
			targetPtr.Words[targetWordNumber], star,
			relation)
	} else {
		fmt.Printf("NIL RELATION\n")
	}
}
