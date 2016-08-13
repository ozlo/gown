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

    exceptions []map[string]string = []map[string]string {
        map[string]string{},    // noun
        map[string]string{},    // verb
        map[string]string{},    // adjective
        map[string]string{},    // adverb
    }
)

func InitiMorphData(dictDirname string) {
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
            exceptions[posIndex][derivedForm] = baseForm
        }
        infile.Close()
    }
}


// Returns a lemmatizations for the word. We assume the word is in the
// raw form. (i.e. spaces are spaces, not '_') This algorithim is similar to,
// but not exactly the same as the Wordnet Morphy algorithm.
func (wn *WN) Morph(origword string, partOfSpeech int) string {
    partOfSpeechIndex := getPosIndex(partOfSpeech)
    if partOfSpeechIndex < 0  {
        // no idea, it's not a supported part of speech
        return origword
    }

    // check the exception lists
    lemma, exists := exceptions[partOfSpeechIndex][origword]
    if exists {
        return lemma
    } else {
        lemma = origword
    }

    if partOfSpeech == POS_ADVERB {
        // only use the exception lists for adverbs
        return origword
    } else {
        // replace the suffxes
        if partOfSpeech == POS_NOUN {
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
                    resp := wn.Lookup(possibleLemma)
                    if len(resp) > 0 {
                        // found it!
                        return possibleLemma
                    }
                }
            }
        }

        return origword
    }
}

func getPosIndex(pos int) int {
    if pos == POS_ADJECTIVE_SATELLITE {
        pos = POS_ADJECTIVE
    }
    return pos - 1
}