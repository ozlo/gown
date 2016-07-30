# GOWN
A Go implementation of the WordNet API.

## Requirements
* WordNet database files https://wordnet.princeton.edu/wordnet/download/current-version/

## WordNet Files
### `cntlist` and `cntlist.rev`
These files contain semantic concordance tagged with the WordNet sense number.
However, these files have not been updated since 2001, and are no longer
maintained, and so these files are ignored.

### `*.pl`
Prolog loadable files are not supported.

### `index.sense`
An index for looking up synsets related to a specific synset.

### Morphology Exception Lists
* `noun.exc`
* `verb.exc`
* `adj.exc`
* `adv.exc`

### WordNet Database Files
* `index.noun`
* `data.noun`
* `index.verb`
* `data.verb`
* `index.adj`
* `data.adj`
* `index.adv`
* `data.adv`

### illustrating Sentences For Verbs
* `sentidx.vrb`
* `sents.vrb`
