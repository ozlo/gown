package gown

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

// syntactic category / part of speech
const POS_UNSUPPORTED int = -1
const POS_NOUN int = 1
const POS_VERB int = 2
const POS_ADJECTIVE int = 3
const POS_ADVERB int = 4
const POS_ADJECTIVE_SATELLITE int = 5

type SenseIndex struct {
    lemmas_to_entries map[string][]SenseIndexEntry
    LEXOGRAPHER_FILE_NUM_TO_NAME []string
}

type SenseIndexEntry struct {
    ss_type int             // POS tag. (e.g. POS_NOUN, ...)
    lex_filenum int         // index into LEXOGRAPHER_FILE_NUM_TO_NAME
    lex_id int              // identifies a sense within a lemma file (default is 0)
    head_word string        // OPTIONAL lemma of the first word of the adjective satellite's head synset. (ss_type of this entry is 5)
    head_id int             // OPTIONAL uniquely identifies head_word in a lexographer file. ( fmt.Sprintf("%s%2d", head_word, head_id) )

    synset_offset int       // byte offset into <POS> data file
    sense_number int        // sense number within the POS for the word
    tag_cnt int             // number of times the word was tagged in semantic concordance texts
}

func loadSenseIndex(senseIndexFilename string) (*SenseIndex, error) {
    index := SenseIndex {
        make(map[string][]SenseIndexEntry),
        []string{
        	"adj.all",            // all adjective clusters
        	"adj.pert",           // relational adjectives (pertainyms)
        	"adv.all",            // all adverbs
        	"noun.Tops",          // unique beginner for nouns
        	"noun.act",           // nouns denoting acts or actions
        	"noun.animal",        // nouns denoting animals
        	"noun.artifact",      // nouns denoting man-made objects
        	"noun.attribute",     // nouns denoting attributes of people and objects
        	"noun.body",          // nouns denoting body parts
        	"noun.cognition",     // nouns denoting cognitive processes and contents
        	"noun.communication", // nouns denoting communicative processes and contents
        	"noun.event",         // nouns denoting natural events
        	"noun.feeling",       // nouns denoting feelings and emotions
        	"noun.food",          // nouns denoting foods and drinks
        	"noun.group",         // nouns denoting groupings of people or objects
        	"noun.location",      // nouns denoting spatial position
        	"noun.motive",        // nouns denoting goals
        	"noun.object",        // nouns denoting natural objects (not man-made)
        	"noun.person",        // nouns denoting people
        	"noun.phenomenon",    // nouns denoting natural phenomena
        	"noun.plant",         // nouns denoting plants
        	"noun.possession",    // nouns denoting possession and transfer of possession
        	"noun.process",       // nouns denoting natural processes
        	"noun.quantity",      // nouns denoting quantities and units of measure
        	"noun.relation",      // nouns denoting relations between people or things or ideas
        	"noun.shape",         // nouns denoting two and three dimensional shapes
        	"noun.state",         // nouns denoting stable states of affairs
        	"noun.substance",     // nouns denoting substances
        	"noun.time",          // nouns denoting time and temporal relations
        	"verb.body",          // verbs of grooming, dressing and bodily care
        	"verb.change",        // verbs of size, temperature change, intensifying, etc.
        	"verb.cognition",     // verbs of thinking, judging, analyzing, doubting
        	"verb.communication", // verbs of telling, asking, ordering, singing
        	"verb.competition",   // verbs of fighting, athletic activities
        	"verb.consumption",   // verbs of eating and drinking
        	"verb.contact",       // verbs of touching, hitting, tying, digging
        	"verb.creation",      // verbs of sewing, baking, painting, performing
        	"verb.emotion",       // verbs of feeling
        	"verb.motion",        // verbs of walking, flying, swimming
        	"verb.perception",    // verbs of seeing, hearing, feeling
        	"verb.possession",    // verbs of buying, selling, owning
        	"verb.social",        // verbs of political and social activities and events
        	"verb.stative",       // verbs of being, having, spatial relations
        	"verb.weather",       // verbs of raining, snowing, thawing, thundering
        	"adj.ppl",            // participial adjectives
        },
    }

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
        lemma := strings.Replace(sense_key_fields[0], "_", " ", -1)
        lex_sense_fields := strings.Split(sense_key_fields[1], ":")
        ss_type, _ := strconv.Atoi(lex_sense_fields[0])     // POS tag. (e.g. POS_NOUN, ...)
        lex_filenum, _ := strconv.Atoi(lex_sense_fields[1]) // index into LEXOGRAPHER_FILE_NUM_TO_NAME
        lex_id, _ := strconv.Atoi(lex_sense_fields[2])      // identifies a sense within a lemma file (default is 0)
        head_word := lex_sense_fields[3]                    // OPTIONAL lemma of the first word of the adjective satellite's head synset. (ss_type of this entry is 5)
        head_id, _ := strconv.Atoi(lex_sense_fields[4])     // OPTIONAL uniquely identifies head_word in a lexographer file. ( fmt.Sprintf("%s%2d", head_word, head_id) )

        newEntry := SenseIndexEntry {
            ss_type,
            lex_filenum,
            lex_id,
            head_word,
            head_id,
            synset_offset,
            sense_number,
            tag_cnt,
        }

        entries, exists := index.lemmas_to_entries[lemma]
        if !exists {
            index.lemmas_to_entries[lemma] = make([]SenseIndexEntry, 1)
            index.lemmas_to_entries[lemma][0] = newEntry
        } else {
            index.lemmas_to_entries[lemma] = append(entries, newEntry)
        }
    }

    return &index, nil
}

 func getLemmaKey(lemma string, sense_id int) string {
     return fmt.Sprintf("%s%02d", lemma, sense_id)
 }

/*
func main() {
    z, _ := loadSenseIndex("./wn-dict/index.sense")
    fmt.Printf("whig = %v\n", z.lemmas_to_entries["whig"])
}
*/
