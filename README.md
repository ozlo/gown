# GOWN
A Go implementation of the WordNet API.

## Requirements
* WordNet database files https://wordnet.princeton.edu/wordnet/download/current-version/

## WordNet Files Utilized
### `index.sense`
An index for looking up synsets related to a specific synset.

### WordNet Database Files
* `index.noun`
* `data.noun`
* `index.verb`
* `data.verb`
* `index.adj`
* `data.adj`
* `index.adv`
* `data.adv`

### Morphology Exception Lists
* `noun.exc`
* `verb.exc`
* `adj.exc`
* `adv.exc`

# TODO
* *Support troponyms for verbs.* This requires adding a inverse relation for all verb hypernyms.
* *Better support for verb groups.* Add fully connect words in a verb groups
