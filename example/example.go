package main

import (
  "fmt"
  "github.com/jonathankoren/gown"
)

func main() {
  wn, err := gown.LoadWordNet("./wn-dict")
  if err != nil {
    fmt.Printf("can't load WordNet: %v\n", err)
    return
  }

  // We'll load the word "live" which can be used as verb, adjective, or adverb
  fmt.Printf("live\n")
  for resultId, senseIndexEntry := range wn.Lookup("live") {
    printSenseIndexEntry(wn, resultId, senseIndexEntry)
    fmt.Printf("\n")
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
  if synsetPtr == nil {
    fmt.Printf("\tNO SYNSET!\n")
  } else {
    //fmt.Printf("\t%v\n", *synsetPtr)
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
        fmt.Printf("\t\t%v\n", relation)
        fmt.Printf("\t\t%d: %s (%d) >> ", i, gown.RELATIONSHIP_ID_TO_STRING[relation.RelationshipType], relation.RelationshipType)
        targetPtr := wn.GetSynset(relation.PartOfSpeech, relation.SynsetOffset)
        if targetPtr != nil {
            fmt.Printf("%d (%d) %v\n", relation.PartOfSpeech, relation.SynsetOffset, targetPtr.Words)
        } else {
            fmt.Printf("NIL RELATION\n")
        }
    }

  }
}
