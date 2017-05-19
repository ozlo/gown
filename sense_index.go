package gown

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

/*
From senseidx(5WN):

The WordNet sense index provides an alternate method for accessing synsets and
word senses in the WordNet database. It is useful to applications that retrieve
synsets or other information related to a specific sense in WordNet, rather
than all the senses of a word or collocation. It can also be used with tools
like grep and Perl to find all senses of a word in one or more parts of speech.
A specific WordNet sense, encoded as a sense_key , can be used as an index into
this file to obtain its WordNet sense number, the database byte offset of the
synset containing the sense, and the number of times it has been tagged in the
semantic concordance texts.
*/

type senseIndex map[string][]SenseIndexEntry
type SenseIndexEntry struct {
    PartOfSpeech int       // POS tag. (e.g. POS_NOUN, ...)
    LexographerFilenum int // index into LEXOGRAPHER_FILE_NUM_TO_NAME
    LexId int              // identifies a sense within a lemma file (default is 0)
    HeadWord string        // OPTIONAL lemma of the first word of the adjective satellite's head synset. (PartOfSpeech of this entry is 5)
    HeadId int             // OPTIONAL uniquely identifies head_word in a lexographer file. ( fmt.Sprintf("%s%2d", head_word, head_id) )

    SynsetOffset int       // byte offset into <POS>.data file
    SenseNumber int        // sense number within the <POS>.data for the word
    TagCount int           // number of times the word was tagged in semantic concordance texts
    synsetPtr *Synset      // back ponter to the underlying synset.
}

func (e *SenseIndexEntry) ToString() string {
    var pos_str string
    switch(e.PartOfSpeech) {
    case POS_UNSUPPORTED:
        pos_str = "UNSUPPORTED"
    case POS_NOUN:
        pos_str = "NOUN"
    case POS_VERB:
        pos_str = "VERB"
    case POS_ADJECTIVE:
        pos_str = "ADJ"
    case POS_ADVERB:
        pos_str = "ADV"
    case POS_ADJECTIVE_SATELLITE:
        pos_str = "ADJ_SAT"
    }

    return fmt.Sprintf("{ %s, file: %s, lex_id: %d head: %s, head_id: %d, synset_offset: %d, sense_number: %d, tag_cnt: %d }",
        pos_str,
        LEXOGRAPHER_FILE_NUM_TO_NAME[e.LexographerFilenum],
        e.LexId,
        e.HeadWord,
        e.HeadId,
        e.SynsetOffset,
        e.SenseNumber,
        e.TagCount)
}

func (e *SenseIndexEntry) GetSynsetPtr() *Synset {
    return e.synsetPtr
}

func loadSenseIndex(wn *WN, senseIndexFilename string) (senseIndex, error) {
    index := senseIndex{}

    infile, err := os.Open(senseIndexFilename)
    defer infile.Close()
    if err != nil {
        return nil, fmt.Errorf("can't open %s: %v", senseIndexFilename, err)
    }
    r := bufio.NewReader(infile)
    if (r == nil) {
        return nil, fmt.Errorf("can't read %s: %v" + senseIndexFilename, err)
    }

    var readerr error = nil
    for ; readerr == nil ; {
        bytebuf, readerr := r.ReadBytes('\n')
        if readerr != nil && readerr != io.EOF {
            panic(readerr)
        }
        if len(bytebuf) == 0 {
            break;
        }

        fields := strings.Split(strings.TrimSpace(string(bytebuf)), " ")
        sense_key := fields[0]
        synset_offset, _ := strconv.Atoi(fields[1])     // byte offset into <POS> data file
        sense_number, _ := strconv.Atoi(fields[2])      // sense number within the POS for the word
        tag_cnt, _ := strconv.Atoi(fields[3])           // number of times the word was tagged in semantic concordance texts

        // now parse the sense key
        sense_key_fields := strings.Split(sense_key, "%")
        lemma := readStoredLemma(sense_key_fields[0])
        lex_sense_fields := strings.Split(sense_key_fields[1], ":")
        ss_type, _ := strconv.Atoi(lex_sense_fields[0])     // POS tag. (e.g. POS_NOUN, ...)
        lex_filenum, _ := strconv.Atoi(lex_sense_fields[1]) // index into LEXOGRAPHER_FILE_NUM_TO_NAME
        lex_id, _ := strconv.Atoi(lex_sense_fields[2])      // identifies a sense within a lemma file (default is 0)
        head_word := lex_sense_fields[3]                    // OPTIONAL lemma of the first word of the adjective satellite's head synset. (ss_type of this entry is 5)
        head_id, _ := strconv.Atoi(lex_sense_fields[4])     // OPTIONAL uniquely identifies head_word in a lexographer file. ( fmt.Sprintf("%s%2d", head_word, head_id) )

        var synsetPtr *Synset = nil
        if wn != nil {
            synsetPtr = wn.GetSynset(ss_type, synset_offset)
        }

        newEntry := SenseIndexEntry {
            ss_type,
            lex_filenum,
            lex_id,
            head_word,
            head_id,
            synset_offset,
            sense_number,
            tag_cnt,
            synsetPtr,
        }

        entries, exists := index[lemma]
        if !exists {
            index[lemma] = make([]SenseIndexEntry, 1)
            index[lemma][0] = newEntry
        } else {
            index[lemma] = append(entries, newEntry)
        }
    }

    return index, nil
}
