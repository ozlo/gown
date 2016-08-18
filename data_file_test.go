package gown

import (
    "testing"
)

var (
    POS_FILE_NAMES = []string { "adj", "adv", "noun", "verb" }
)

func TestLoadPosIndex(t *testing.T) {
    dictDir, _ := GetWordNetDictDir()
    for _, posName := range POS_FILE_NAMES {
        posIndexFilename := dictDir + "/index."  + posName
        _, err := readPosIndex(posIndexFilename)
        if err != nil {
            t.Fatalf("failed to read %s: %v", posIndexFilename, err)
        }
    }
}

func TestLoadPosData(t *testing.T) {
    dictDir, _ := GetWordNetDictDir()
    for _, posName := range POS_FILE_NAMES {
        posDataFilename := dictDir + "/data."  + posName
        _, err := readPosData(posDataFilename)
        if err != nil {
            t.Fatalf("failed to read %s: %v", posDataFilename, err)
        }
    }
}
