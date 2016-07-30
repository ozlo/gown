package gown

type WN struct {
    senseIndex *senseIndex
    posIndicies map[int]*dataIndex
    posData map[int]*dataFile
}

func LoadWordNet(dictDirname string) (*WN, error) {
    wn := WN {
        senseIndex: nil,
        posIndicies: map[int]*dataIndex{},
        posData: map[int]*dataFile{},
    }

    var err error = nil
    wn.senseIndex, err = loadSenseIndex(dictDirname + "/index.sense")
    if err != nil {
        return nil, err
    }

    pos_file_names := []string { "", "noun", "verb", "adj", "adv" }
    for i := 1; i < len(pos_file_names); i++ {
        wn.posIndicies[i], err = readPosIndex(dictDirname + "/index." + pos_file_names[i])
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
    posIndexPtr, exists := wn.posIndicies[pos]
    if !exists {
        return nil
    }
    sn, exists := (*posIndexPtr)[lemma]
    if exists {
        return &sn
    } else {
        return nil
    }
}

func (wn *WN) Lookup(lemma string) []*SenseIndexEntry {
    senseEntries, exists := (*wn.senseIndex)[lemma]
    if !exists {
        return []*SenseIndexEntry{}
    }
    ret := make([]*SenseIndexEntry, len(senseEntries))
    for i, e := range senseEntries {
        ret[i] = &e
    }
    return ret
}
