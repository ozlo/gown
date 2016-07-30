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

  fmt.Printf("\n===================\n\n")

  // now let's lookup a word with a known part of speech
  dataIndexEntry := wn.LookupWithPartOfSpeech("computer", gown.POS_NOUN)
  if dataIndexEntry == nil {
      fmt.Printf("Can't found a computer as a noun!\n")
  } else {
      fmt.Printf("computer (noun)\n")
      fmt.Printf("%v\n", *dataIndexEntry)
      for _, synsetOffset := range dataIndexEntry.SynsetOffsets {
          printSynsetPtr(wn, wn.GetSynset(gown.POS_NOUN, synsetOffset))
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
          printRelationship(wn, i, relation, synsetPtr.Words)
      }
    }
}

func printRelationship(wn *gown.WN, i int, relation gown.RelationshipEdge, srcWords []string) {
    //fmt.Printf("\t\t%v\n", relation)
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
