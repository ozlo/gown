package gown

import (
    "testing"
)

func BenchmarkLoadWordNet(b *testing.B) {
    dictDir, _ := GetWordNetDictDir()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        LoadWordNet(dictDir)
    }
}

func BenchmarkLookupWithPartOfSpeech(b *testing.B) {
    dictDir, _ := GetWordNetDictDir()
    wn, _ := LoadWordNet(dictDir)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wn.LookupWithPartOfSpeech("computer", POS_NOUN)
    }
}

func BenchmarkLookup(b *testing.B) {
    dictDir, _ := GetWordNetDictDir()
    wn, _ := LoadWordNet(dictDir)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wn.Lookup("live")
    }
}

func TestMorph(t *testing.T) {
    dictDir, _ := GetWordNetDictDir()
    wn, _ := LoadWordNet(dictDir)
    wn.InitMorphData(dictDir)
    poses := []int {
        POS_VERB, // are
        POS_NOUN, // splits
        POS_VERB, // splits
        POS_VERB, // left
        POS_NOUN, // trucks
        POS_VERB, // saw
        POS_NOUN, // children
        POS_VERB, // swam
        POS_NOUN, // remains
        POS_VERB, // remains
        POS_NOUN, // plant
        POS_NOUN, // Angus
        POS_NOUN, // Angus
        POS_VERB, // walked
        POS_NOUN, // park
        POS_VERB, // jumping
        POS_NOUN, // octopuses
        POS_NOUN, // octopi
        POS_NOUN, // octopus
        POS_NOUN, // ewoks
        POS_VERB, // wanted
        POS_VERB, // grilled
    }
    inputs := []string {
        "are",
        "splits",
        "splits",
        "left",
        "trucks",
        "saw",
        "children",
        "swam",
        "remains",
        "remains",
        "plant",
        "Angus",
        "angus",
        "walked",
        "park",
        "jumping",
        "octopuses",
        "octopi",
        "octopus",
        "ewoks",
        "wanted",
        "grilled",
    }
    expecteds := []string {
        "be",
        "split",
        "split",
        "leave",
        "truck",
        "see",
        "child",
        "swim",
        "remains",
        "remain",
        "plant",
        "Angus",
        "angus",
        "walk",
        "park",
        "jump",
        "octopus",
        "octopus",
        "octopus",
        "", // ewok
        "want",
        "grill",
    }

    for i, pos := range poses {
        input := inputs[i]
        expected := expecteds[i]
        actual := wn.Morph(input, pos)
        if actual != expected {
            t.Errorf("for %s/%d expected %s but got %s\n", input, pos, expected, actual)
        }
    }
}
