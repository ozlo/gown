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
    partOfSpeech int
    synsetCount int
    relationships []int
    tagSenseCount int
    synsetOffsets []int
}
func (die *DataIndexEntry) GetPartOfSpeech() int {
    return die.partOfSpeech
}
func (die *DataIndexEntry) GetSynsetCount() int {
    return die.synsetCount
}
func (die *DataIndexEntry) GetRelationships() []int {
    return copyIntArray(die.relationships)
}
func (die *DataIndexEntry) GetTagSenseCount() int {
    return die.tagSenseCount
}
func (die *DataIndexEntry) GetSynsetOffsets() []int {
    return copyIntArray(die.synsetOffsets)
}

type dataFile map[int]*Synset
type Synset struct {
    synsetOffset int
    lexographerFilenum int
    partOfSpeech int
    words []string
    lexIds []int
    relationships []RelationshipEdge
    gloss string
}
func (s *Synset) GetSynsetOffset() int {
    return s.synsetOffset
}
func (s *Synset) GetLexographerFilenum() int {
    return s.lexographerFilenum
}
func (s *Synset) GetLexographerFilename() string {
    return cLexographerFileNumToName[s.lexographerFilenum]
}
func (s *Synset) GetPartOfSpeech() int {
    return s.partOfSpeech
}
func (s *Synset) GetWords() []string {
    return copyStringArray(s.words)
}
func (s *Synset) GetLexIds() []int {
    return copyIntArray(s.lexIds)
}
func (s *Synset) GetRelationships() []RelationshipEdge {
    ret := make([]RelationshipEdge, len(s.relationships))
    for i, e := range s.relationships {
        ret[i] = e
    }
    return ret
}
func (s *Synset) GetGloss() string {
    return s.gloss
}

type RelationshipEdge struct {
    relationshipType int      // ANTONYM_RELATIONSHIP, etc.
    synsetOffset int          // synset offset of the target
    partOfSpeech int          // part-of-speech of target
    sourceWordNumber int      // word number of the source
    targetWordNumber int      // word number of the target
}
func (re *RelationshipEdge) GetRelationshipType() int {
    return re.relationshipType
}
func (re *RelationshipEdge) GetSynsetOffset() int {
    return re.synsetOffset
}
func (re *RelationshipEdge) GetPartOfSpeech() int {
    return re.partOfSpeech
}
func (re *RelationshipEdge) GetSourceWordNumber() int {
    return re.sourceWordNumber
}
func (re *RelationshipEdge) GetTargetWordNumber() int {
    return re.targetWordNumber
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
            relationships[i], _ = cRelationshipPointerSymbols[fields[field_index]]
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
            partOfSpeech: pos_tag,
            synsetCount: synset_cnt,
            relationships: relationships,
            tagSenseCount: tagsense_cnt,
            synsetOffsets: synsetOffsets,
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
            pointer_type, symbolFound := cRelationshipPointerSymbols[fields[fieldIndex]]
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
                relationshipType: pointer_type,
                synsetOffset: synset_offset,
                partOfSpeech: pos,
                sourceWordNumber: src_word_num,
                targetWordNumber: dest_word_num,
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
                synsetOffset: synset_offset,
                lexographerFilenum: lex_filenum,
                partOfSpeech: ss_type,
                words: words,
                lexIds: lex_ids,
                relationships: pointers,
                gloss: gloss,
        }
    }

    return data, nil
}
