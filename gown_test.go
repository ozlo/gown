package gown

import (
    "testing"
)

func BenchmarkLoadWordNet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        LoadWordNet("./wn-dict")
    }
}

func BenchmarkLookupWithPartOfSpeech(b *testing.B) {
    wn, _ := LoadWordNet("./wn-dict")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wn.LookupWithPartOfSpeech("computer", POS_NOUN)
    }
}

func BenchmarkLookup(b *testing.B) {
    wn, _ := LoadWordNet("./wn-dict")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wn.Lookup("live")
    }
}
