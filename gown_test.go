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
