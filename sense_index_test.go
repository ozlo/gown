package gown

import (
    "testing"
)

func TestLoadSenseIndex(t *testing.T) {
    dictDir, _ := GetWordNetDictDir()
    senseIndexFile := dictDir + "/index.sense"
    senseIndex, err := loadSenseIndex(nil, senseIndexFile)
    if senseIndex == nil {
        t.Fatalf("Failed to load sense index: %v", err)
    }

    // ------------------------------------------------------------------------
    // validate "computer"
    /*
    computer: {1  6 0 0 3086983 1 6} { NOUN, file: noun.artifact, lex_id: 0 head: , head_id: 0, synset_offset: 3086983, sense_number: 1, tag_cnt: 6 }
    computer: {1 18 0 0 9906486 2 0} { NOUN, file: noun.person,   lex_id: 0 head: , head_id: 0, synset_offset: 9906486, sense_number: 2, tag_cnt: 0 }
    */
    computerLemmas, _ := senseIndex["computer"]
    if computerLemmas == nil || len(computerLemmas) == 0 {
        t.Fatalf("\"computer\" not found in sense index. Not loaded correctly?")
    }
    if len(computerLemmas) != 2 {
        t.Fatalf("expected 2 lemmas for \"computer\", but got %d\n", len(computerLemmas))
    }
    expected_computer_pos := map[int]int { POS_NOUN: 2 }
    actual_computer_pos := map[int]int {}
    expected_computer_files := map[string]int {
        "noun.artifact": 1,
        "noun.person": 1,
    }
    actual_computer_files := map[string]int {}
    for _, lemma_index := range computerLemmas {
        t.Logf("%s: %v %s\n", "computer", lemma_index, lemma_index.ToString())
        actual_computer_pos[lemma_index.GetPartOfSpeech()]++
        actual_computer_files[lemma_index.GetLexographerFilename()]++
    }
    validate_pos(t, "computer", expected_computer_pos, actual_computer_pos)
    validate_counts(t, "computer", expected_computer_files, actual_computer_files)
    t.Logf("\n")


    // ------------------------------------------------------------------------
    // validate "live"
    /*
    live%2:31:00:: 00598039 6 1
    live%2:42:00:: 02620216 5 14
    live%2:42:01:: 02622766 4 16
    live%2:42:04:: 02624202 3 29
    live%2:42:06:: 02620422 2 51
    live%2:42:07:: 02621023 7 0
    live%2:42:08:: 02655932 1 129
    live%3:00:00:: 00095301 3 0
    live%3:00:01:: 00100143 2 0
    live%3:00:02:: 01425429 1 6
    live%4:02:00:: 00260451 1 0
    live%5:00:00:charged:00 00359472 10 0
    live%5:00:00:current:00 00670576 9 0
    live%5:00:00:elastic:00 00847134 6 0
    live%5:00:00:lively:00 00809813 7 0
    live%5:00:00:loaded:00 01427469 5 0
    live%5:00:00:reverberant:00 02017537 4 0
    live%5:00:02:current:00 00670686 8 0
    live%5:00:07:active:05 00041710 11 0
    */
    liveLemmas, _ := senseIndex["live"]
    if liveLemmas == nil || len(liveLemmas) == 0 {
        t.Fatalf("\"live\" not found in sense index. Not loaded correctly?")
    }
    if len(liveLemmas) != 19 {
        t.Fatalf("expected 19 lemmas for \"live\", but got %d\n", len(liveLemmas))
    }
    expected_live_pos := map[int]int {
        POS_VERB: 7,
        POS_ADJECTIVE: 3,
        POS_ADVERB: 1,
        POS_ADJECTIVE_SATELLITE: 8,
    }
    actual_live_pos := map[int]int {}
    expected_live_files := map[string]int {
        "verb.cognition": 1,
        "verb.stative": 6,
        "adj.all": 11,
        "adv.all": 1,
    }
    actual_live_files := map[string]int {}
    for _, lemma_index := range liveLemmas {
        t.Logf("%s: %v %s\n", "live", lemma_index, lemma_index.ToString())
        actual_live_pos[lemma_index.GetPartOfSpeech()]++
        actual_live_files[lemma_index.GetLexographerFilename()]++
    }
    validate_pos(t, "live", expected_live_pos, actual_live_pos)
    validate_counts(t, "live", expected_live_files, actual_live_files)
    t.Logf("\n")

    // ------------------------------------------------------------------------
    // validate "Angus"
    /*
    02408581 05 n 03 Aberdeen_Angus 0 Angus 0 black_Angus 0 001 @ 02406838 n 0000 | black hornless breed from Scotland
    */
    angusLemmas, _ := senseIndex["angus"]
    t.Logf("angusLemmas = %v\n", angusLemmas)
    if angusLemmas == nil || len(angusLemmas) == 0 {
        t.Fatalf("\"angus\" not found in sense index. Not loaded correctly?")
    }
    if len(angusLemmas) != 2 {
        t.Fatalf("expected 2 lemmas for \"angus\", but got %d\n", len(angusLemmas))
    }
    expected_angus_pos := map[int]int { POS_NOUN: 2 }
    actual_angus_pos := map[int]int {}
    expected_angus_files := map[string]int {
        "noun.animal": 1,
        "noun.person": 1,
    }
    actual_angus_files := map[string]int {}
    for _, lemma_index := range angusLemmas {
        t.Logf("%s: %v %s\n", "angus", lemma_index, lemma_index.ToString())
        actual_angus_pos[lemma_index.GetPartOfSpeech()]++
        actual_angus_files[lemma_index.GetLexographerFilename()]++
    }
    validate_pos(t, "angus", expected_angus_pos, actual_angus_pos)
    validate_counts(t, "angus", expected_angus_files, actual_angus_files)
    t.Logf("\n")
}

func validate_pos(t *testing.T, word string, expected map[int]int, actual map[int]int) {
    for k, v := range actual {
        if expected[k] != v {
            t.Fatalf("%v expected pos %v with count %v, but had %v", word, k, expected[k], v)
        }
    }
    for k, v := range expected {
        if actual[k] != v {
            t.Fatalf("%v expected pos %v with count %v, but had %v", word, k, actual[k], v)
        }
    }
    for k, v := range actual {
        if expected[k] != v {
            t.Fatalf("%v expected pos %v with count %v, but had %v", word, k, expected[k], v)
        }
    }
    for k, v := range expected {
        if actual[k] != v {
            t.Fatalf("%v expected pos %v with count %v, but had %v", word, k, actual[k], v)
        }
    }
}

func validate_counts(t *testing.T, word string, expected map[string]int, actual map[string]int) {
    for k, v := range actual {
        if expected[k] != v {
            t.Fatalf("%v expected file %v with count %v, but had %v", word, k, expected[k], v)
        }
    }
    for k, v := range expected {
        if actual[k] != v {
            t.Fatalf("%v expected file %v with count %v, but had %v", word, k, actual[k], v)
        }
    }
    for k, v := range actual {
        if expected[k] != v {
            t.Fatalf("%v expected file %v with count %v, but had %v", word, k, expected[k], v)
        }
    }
    for k, v := range expected {
        if actual[k] != v {
            t.Fatalf("%v expected file %v with count %v, but had %v", word, k, actual[k], v)
        }
    }
}
