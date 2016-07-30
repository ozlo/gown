package gown

import (
    "testing"
)

func TestLoadPosIndex(t *testing.T) {
    for _, posName := range POS_FILE_NAMES {
        posIndexFilename := "./wn-dict/index."  + posName
        _, err := readPosIndex(posIndexFilename)
        if err != nil {
            t.Errorf("failed to read %s: %v", posIndexFilename, err)
        }
    }
}

func TestLoadPosData(t *testing.T) {
    for _, posName := range POS_FILE_NAMES {
        posDataFilename := "./wn-dict/data."  + posName
        _, err := readPosData(posDataFilename)
        if err != nil {
            t.Errorf("failed to read %s: %v", posDataFilename, err)
        }
    }
}
