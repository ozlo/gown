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

    addr = [][]string {
        // noun
        []string { "", "s", "x", "z", "ch", "sh", "man", "y" },
        // verb
        []string { "", "y", "e", "", "e", "", "e", "" },
        // adjective
        []string { "", "", "e", "e" },
    }

    prepositions = map[string]struct{} {
        "to": struct{}{},
        "at": struct{}{},
        "of": struct{}{},
        "on": struct{}{},
        "off": struct{}{},
        "in": struct{}{},
        "out": struct{}{},
        "up": struct{}{},
        "down": struct{}{},
        "from": struct{}{},
        "with": struct{}{},
        "into": struct{}{},
        "for": struct{}{},
        "about": struct{}{},
        "between": struct{}{},
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


// Returns a set of lemmatizations for the word. We assume the word is in the
// raw form. (i.e. spaces are spaces, not '_')
func Morph(origword string, partOfSpeech int) []string {
    partOfSpeechIndex := getPosIndex(partOfSpeech)
    if partOfSpeechIndex < 0  {
        // no idea, it's not a supported part of speech
        return []string { origword }
    }

    // check the exception lists
    lemma, exists := exceptions[partOfSpeechIndex][origword]
    if exists {
        return []string{ lemma }
    }

    if partOfSpeech == POS_VERB {
        toks := strings.Split(" ", origword)
        ending := ""
        if len(toks) > 1 {
            ending = " " + strings.Join(toks[1:], " ")
        }
        results := []string{}
        for _, lemma := range findLemmas(toks[0], partOfSpeech) {
            results = append(results, lemma + ending)
        }
        return results
    } else {
        return morphword(origword, partOfSpeech)
    }
}

func morphword(origword string, partOfSpeech int) []string {
    lemmas := []string{}
    partOfSpeechIndex := getPosIndex(partOfSpeech)
    exception, exists := exceptions[partOfSpeechIndex][origword]
    if exists {
        lemmas = append(lemmas, exception)
    }

    if partOfSpeech == POS_ADVERB {
         // skip it
        if len(lemmas) == 0 {
            return []string { origword }
        } else {
            return lemmas
        }
    } else if partOfSpeech == POS_NOUN {
        if strings.HasSuffix(origword, "ful") {
            origword = origword[:len(origword) - 3]
        } else {
            if strings.HasSuffix(origword, "ss") || len(origword) <= 2 {
                return []string { origword }
            }
        }
    }

    return append(lemmas, findLemmas(origword, partOfSpeech)...)
}

func getPosIndex(partOfSpeech int) int {
    if partOfSpeech == POS_ADJECTIVE_SATELLITE {
        partOfSpeech = POS_ADJECTIVE
    }
    return partOfSpeech - 1
}

func findLemmas(origword string, partOfSpeech int) []string {
    partOfSpeechIndex := getPosIndex(partOfSpeech)
    exception, exists := exceptions[partOfSpeechIndex][origword]
    if exists {
        return []string { exception }
    } else {
        lemmas := []string{}
        wordlen := len(origword)
        for lastIndex := 1; lastIndex <= 4; lastIndex++ {
            if (wordlen - lastIndex) > 0 {
                suffix := origword[lastIndex:]
                replacements, suffixFound := suffixReplacements[partOfSpeechIndex][suffix]
                if suffixFound {
                    prefix := origword[:lastIndex]
                    for _, replace := range replacements {
                        lemmas = append(lemmas, prefix + replace)
                    }
                }
            } else {
                // not enough letters
                break
            }
        }

        if len(lemmas) == 0 {
            return []string { origword }
        } else {
            return lemmas
        }
    }
}
