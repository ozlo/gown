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
From wndb(5WN):
For each syntactic category, two files are needed to represent the contents of
the WordNet database - index. pos and data. pos, where pos is noun, verb,
adj and adv . The other auxiliary files are used by the WordNet library's
searching functions and are needed to run the various WordNet browsers.

Each index file is an alphabetized list of all the words found in WordNet in
the corresponding part of speech. On each line, following the word, is a list
of byte offsets (synset_offset s) in the corresponding data file, one for each
synset containing the word. Words in the index file are in lower case only,
regardless of how they were entered in the lexicographer files. This folds
various orthographic representations of the word into one line enabling
database searches to be case insensitive. See wninput(5WN) for a detailed
description of the lexicographer files

A data file for a syntactic category contains information corresponding to the
synsets that were specified in the lexicographer files, with relational
pointers resolved to synset_offset s. Each line corresponds to a synset.
Pointers are followed and hierarchies traversed by moving from one synset to
another via the synset_offset s.
*/

type dataIndex map[string]DataIndexEntry
type DataIndexEntry struct {
    PartOfSpeech int
    SynsetCount int
    Relationships []int
    TagSenseCount int
    SynsetOffsets []int
}

type dataFile map[int]*Synset
type Synset struct {
    SynsetOffset int
    LexographerFilenum int
    PartOfSpeech int
    Words []string
    LexIds []int
    Relationships []RelationshipEdge
    Gloss string
}
type RelationshipEdge struct {
    RelationshipType int      // ANTONYM_RELATIONSHIP, etc.
    SynsetOffset int          // synset offset of the target
    PartOfSpeech int          // part-of-speech of target
    SourceWordNumber int      // word number of the source
    TargetWordNumber int      // word number of the target
}

type DataIndexPair struct {
    Lexeme string
    IndexEntry DataIndexEntry
}
func DataIndexIterator(di *dataIndex) <-chan DataIndexPair {
    ch := make(chan DataIndexPair)
    go func() {
        for k, v := range *di {
            ch <- DataIndexPair {
                Lexeme: k,
                IndexEntry: v,
            }
        }
        close(ch) // Remember to close or the loop never ends!
    }()
    return ch
}

// Reads a index.POS (e.g. index.noun, index.verb, etc.) file and populates
// a dataIndex . The index format is:
// lemma  pos  synset_cnt  p_cnt  [ptr_symbol...]  sense_cnt  tagsense_cnt   synset_offset  [synset_offset...]
func readPosIndex(posIndexFilename string) (dataIndex, error) {
    index := dataIndex{}

    infile, err := os.Open(posIndexFilename)
    defer infile.Close()
    if err != nil {
        return nil, fmt.Errorf("can't open %s: %v", posIndexFilename, err)
    }
    r := bufio.NewReader(infile)
    if (r == nil) {
        return nil, fmt.Errorf("can't read %s: %v" + posIndexFilename, err)
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
        line := string(bytebuf)
        if line[0:2] == "  " {
            // comment line
            continue
        }
        fields := strings.SplitN(strings.TrimSpace(line), " ", -1)
        lemma := readStoredLemma(fields[0])
        pos_tag := oneCharPosTagToPosId(fields[1])
        synset_cnt, _ := strconv.Atoi(fields[2])     // number of senses of the <lemma, pos> pair
        p_cnt, _ := strconv.Atoi(fields[3])          // number of different pointers that lemma has in all synsets containing it.
        field_index := 4
        relationships := make([]int, p_cnt)
        // consume p_cnt pointer symbols
        for i := 0; i < p_cnt; i++ {
            relationships[i], _ = RELATIONSHIP_POINTER_SYMBOLS[fields[field_index]]
            field_index++
        }
        field_index++  // sense_cnt is redundant with synset_cnt, so skip it
        tagsense_cnt, _ := strconv.Atoi(fields[field_index])
        field_index++
        synsetOffsets := make([]int, synset_cnt)
        for i := 0; i < synset_cnt; i++ {
            synsetOffsets[i], _ = strconv.Atoi(fields[field_index])
            field_index++
        }

        _, exists := index[lemma]
        if exists {
            fmt.Printf("WARNING: %s already exists. Overwriting.\n", lemma)
        }
        index[lemma] = DataIndexEntry {
            PartOfSpeech: pos_tag,
            SynsetCount: synset_cnt,
            Relationships: relationships,
            TagSenseCount: tagsense_cnt,
            SynsetOffsets: synsetOffsets,
        }
    }

    return index, nil
}

// Reads a data.POS (e.g. data.noun, data.verb, etc.) file and populates
// a map of ints to dataIndexEntries. The data format is:
// synset_offset  lex_filenum  ss_type  w_cnt  word  lex_id  [word  lex_id...]  p_cnt  [ptr...]  [frames...]  |   gloss
func readPosData(posDataFilename string) (dataFile, error) {
    data := dataFile{}

    infile, err := os.Open(posDataFilename)
    defer infile.Close()
    if err != nil {
        return nil, fmt.Errorf("can't open %s: %v", posDataFilename, err)
    }
    r := bufio.NewReader(infile)
    if (r == nil) {
        return nil, fmt.Errorf("can't read %s: %v" + posDataFilename, err)
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
        line := string(bytebuf)
        if line[0:2] == "  " {
            // comment line
            continue
        }
        fields := strings.SplitN(strings.TrimSpace(line), " ", -1)
        synset_offset, _ := strconv.Atoi(fields[0])
        lex_filenum, _ := strconv.Atoi(fields[1])
        ss_type := oneCharPosTagToPosId(fields[2])
        w_cnt64, _ := strconv.ParseInt(fields[3], 16, 0)
        w_cnt := int(w_cnt64)
        words := make([]string, w_cnt)
        lex_ids := make([]int, w_cnt)
        fieldIndex := 4
        for i := 0; i < w_cnt; i++ {
            words[i] = readStoredLemma(fields[fieldIndex])
            fieldIndex++
            lex_id64, _ := strconv.ParseInt(fields[fieldIndex], 16, 0)
            lex_ids[i] = int(lex_id64)
            fieldIndex++
        }
        p_cnt, _ := strconv.Atoi(fields[fieldIndex])
        fieldIndex++
        pointers := make([]RelationshipEdge, p_cnt)
        for i := 0; i < p_cnt; i++ {
            pointer_type, symbolFound := RELATIONSHIP_POINTER_SYMBOLS[fields[fieldIndex]]
            if !symbolFound {
                panic(fmt.Sprintf("could not handle relationship symbol %s in line <<%v>>, file %s", fields[fieldIndex], line, posDataFilename))
            }
            fieldIndex++
            synset_offset, _ := strconv.Atoi(fields[fieldIndex])
            fieldIndex++
            pos := oneCharPosTagToPosId(fields[fieldIndex])
            fieldIndex++

            src_wordnum64, _ := strconv.ParseInt(fields[fieldIndex][0:2], 16, 0)
            dest_wordnum64, _ := strconv.ParseInt(fields[fieldIndex][2:4], 16, 0)
            fieldIndex++
            src_word_num := int(src_wordnum64)
            dest_word_num := int(dest_wordnum64)
            pointers[i] = RelationshipEdge {
                RelationshipType: pointer_type,
                SynsetOffset: synset_offset,
                PartOfSpeech: pos,
                SourceWordNumber: src_word_num,
                TargetWordNumber: dest_word_num,
            }
        }
        // skip data.verb frames

        pipeIndex := strings.LastIndex(line, "|")
        var gloss string
        if pipeIndex >= 0 {
            gloss = strings.TrimSpace(line[pipeIndex + 2:])
        } else {
            gloss = ""
        }

        data[synset_offset] = &Synset {
                SynsetOffset: synset_offset,
                LexographerFilenum: lex_filenum,
                PartOfSpeech: ss_type,
                Words: words,
                LexIds: lex_ids,
                Relationships: pointers,
                Gloss: gloss,
        }
    }

    return data, nil
}
