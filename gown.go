package gown

import (
    "fmt"
    "os"
    "strings"
)

type WN struct {
    senseIndex *senseIndex
    PosIndicies map[int]*dataIndex
    posData map[int]*dataFile
    exceptions []map[string]string
}

func GetWordNetDictDir() (string, error) {
    systemDefaults := []string {
        "/usr/WordNet-3.%d/dict",
        "/usr/share/WordNet-3.%d/dict",
        "/usr/local/WordNet-3.%d/dict",
        "/usr/local/share/WordNet-3.%d/dict",
        "/opt/WordNet-3.%d/dict",
        "/opt/share/WordNet-3.%d/dict",
        "/opt/local/WordNet-3.%d/dict",
        "/opt/local/share/WordNet-3.%d/dict",
    }
    // check environment variables
    dictname := os.Getenv("WNHOME") + "/dict"
    _, err := os.Stat(dictname)
    if err == nil {
        return dictname, nil
    }

    dictname = os.Getenv("WNSEARCHDIR")
    _, err = os.Stat(dictname)
    if err == nil {
        return dictname, nil
    }

    // check possible installation dirs
    for v := 0; v <= 1; v++ {   // checks for WordNet 3.0 and 3.1
        for _, systemDefault := range systemDefaults {
            dictname = fmt.Sprintf(systemDefault, v)
            _, err = os.Stat(dictname)
            if err == nil {
                return dictname, nil
            }
        }
    }

    // tried everything
    return "", fmt.Errorf("Can't find WordNet dictionary")
}


func LoadWordNet(dictDirname string) (*WN, error) {
    wn := WN {
        senseIndex: nil,
        PosIndicies: map[int]*dataIndex{},
        posData: map[int]*dataFile{},
    }

    var err error = nil
    wn.senseIndex, err = loadSenseIndex(dictDirname + "/index.sense")
    if err != nil {
        return nil, err
    }

    pos_file_names := []string { "", "noun", "verb", "adj", "adv" }
    for i := 1; i < len(pos_file_names); i++ {
        wn.PosIndicies[i], err = readPosIndex(dictDirname + "/index." + pos_file_names[i])
        if err != nil {
            return nil, err
        }
        wn.posData[i], err = readPosData(dictDirname + "/data." + pos_file_names[i])
        if err != nil {
            return nil, err
        }
    }

    return &wn, nil
}

func (wn *WN) LookupWithPartOfSpeech(lemma string, pos int) *DataIndexEntry {
    posIndexPtr, exists := wn.PosIndicies[pos]
    if !exists {
        return nil
    }
    sn, exists := (*posIndexPtr)[strings.ToLower(lemma)]
    if exists {
        return &sn
    } else {
        return nil
    }
}

func (wn *WN) Lookup(lemma string) []*SenseIndexEntry {
    senseEntries, exists := (*wn.senseIndex)[strings.ToLower(lemma)]
    if !exists {
        return []*SenseIndexEntry{}
    }
    ret := make([]*SenseIndexEntry, len(senseEntries))
    for i, _ := range senseEntries {
        ret[i] = &senseEntries[i]
    }
    return ret
}

func (wn *WN) GetSynset(pos int, synsetOffset int) *Synset {
  if pos == POS_ADJECTIVE_SATELLITE {
    pos = POS_ADJECTIVE
  }
  idxPtr, exists := wn.posData[pos]
  if !exists || idxPtr == nil {
    return nil
  }
  s, exists := (*idxPtr)[synsetOffset]
  if !exists {
    return nil
  }
  return &s
}

func (wn *WN) Iter() <-chan Synset {
    outChan := make(chan Synset)
    go func () {
        for _, datFile := range wn.posData {
            for _, synset := range *datFile {
                words := make([]string, len(synset.Words))
                for i, w := range synset.Words {
                    words[i] = w
                }
                lexids := make([]int, len(synset.LexIds))
                for i, w := range synset.LexIds {
                    lexids[i] = w
                }
                edges := make([]RelationshipEdge, len(synset.Relationships))
                for i, w := range synset.Relationships {
                    edges[i] = w
                }
                outChan <- Synset {
                    SynsetOffset: synset.SynsetOffset,
                    LexographerFilenum: synset.LexographerFilenum,
                    PartOfSpeech: synset.PartOfSpeech,
                    Words: words,
                    LexIds: lexids,
                    Relationships: edges,
                    Gloss: synset.Gloss,
                }
            }
        }
        close(outChan)
    } ();
    return outChan
}
