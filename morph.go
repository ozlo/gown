package gown

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

var (
    suffixReplacements = []map[string][]string {
        // noun
        map[string][]string {
            "s": []string { "" },
            "ses": []string { "s" },
            "xes": []string { "x" },
            "zes": []string { "z" },
            "ches": []string { "ch" },
            "shes": []string { "sh" },
            "men": []string { "man" },
            "ies": []string { "y" },
        },
        // verb
        map[string][]string {
            "s": []string { "" },
            "ies": []string { "y" },
            "es": []string { "e", "" },
            "ed": []string { "e", "" },
            "ing": []string { "e", "" },

        },
        // adjective
        map[string][]string {
            "er": []string { "", "e" },
            "est": []string { "", "e"},
        },
    }
)

func (wn *WN) InitMorphData(dictDirname string) {
    wn.exceptions = []map[string]string {
        map[string]string{},    // noun
        map[string]string{},    // verb
        map[string]string{},    // adjective
        map[string]string{},    // adverb
    }

    posNames := []string { "noun", "verb", "adj", "adv" }
    for posIndex, posName := range posNames {
        exceptionFilename := dictDirname + string(filepath.Separator) + posName + ".exc"
        infile, err := os.Open(exceptionFilename)
        if err != nil {
            panic(fmt.Sprintf("Can't open morph exception file %s: %v", exceptionFilename, err))
        }
        r := bufio.NewReader(infile)
        if (r == nil) {
            panic(fmt.Sprintf("Can't open morph exception file %s: %v" + exceptionFilename, err))
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
            fields := strings.SplitN(strings.TrimSpace(string(bytebuf)), " ", -1)
            derivedForm := strings.Replace(fields[0], "_", " ", -1)
            baseForm := strings.Replace(fields[1], "_", " ", -1)
            wn.exceptions[posIndex][derivedForm] = baseForm
        }
        infile.Close()
    }
}


// Returns a lemmatizations for the word. We assume the word is in the
// raw form. (i.e. spaces are spaces, not '_') This algorithim is similar to,
// but not exactly the same as the Wordnet Morphy algorithm.
//
// This function does not handle prepositional verb phrases,
// If no base morph is found, assumes the original word is the base.
func (wn *WN) Morph(origword string, partOfSpeech int) string {
    partOfSpeechIndex := getPosIndex(partOfSpeech)
    if partOfSpeechIndex < 0  {
        // no idea, it's not a supported part of speech
        return ""
    }

    // check the exception lists
    lemma, exists := wn.exceptions[partOfSpeechIndex][origword]
    if exists {
        return lemma
    }

    if partOfSpeech == POS_ADVERB {
        // only use the exception lists for adverbs
        return ""
    }

    if partOfSpeech != POS_VERB {
        // check the original
        dataIndexEntry := wn.LookupWithPartOfSpeech(origword, partOfSpeech)
        if dataIndexEntry != nil {
            return origword
        }
    }

    if partOfSpeech == POS_NOUN {
        // if it's a noun, drop the -full or -ss suffixes
        if strings.HasSuffix(origword, "ful") {
            origword = origword[:len(origword) - 3]
        } else {
            if strings.HasSuffix(origword, "ss") || len(origword) <= 2 {
                // too small
                return origword
            }
        }
    }

    for i := 1; i <= 4; i++ {
        suffixIndex := len(origword) - i
        if suffixIndex <= 0 {
            break;
        }

        baseword := origword[:suffixIndex]
        suffix := origword[suffixIndex:]
        replacements, found := suffixReplacements[partOfSpeechIndex][suffix]
        if found {
            for _, replacement := range replacements {
                possibleLemma := baseword + replacement
                dataIndexEntry := wn.LookupWithPartOfSpeech(possibleLemma, partOfSpeech)
                if dataIndexEntry != nil {
                    // found it!
                    return possibleLemma
                }
            }
        }
    }

    // failed
    return ""
}

func getPosIndex(pos int) int {
    if pos == POS_ADJECTIVE_SATELLITE {
        pos = POS_ADJECTIVE
    }
    return pos - 1
}
