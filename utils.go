package gown

import (
    "fmt"
    "strings"
)

//var LEXOGRAPHER_FILE_NUM_TO_NAME []string
//var RELATIONSHIP_POINTER_SYMBOLS map[string]int

var (
    POS_FILE_NAMES = []string { "adj", "adv", "noun", "verb" }

    LEXOGRAPHER_FILE_NUM_TO_NAME = []string{
        "adj.all",            // all adjective clusters
        "adj.pert",           // relational adjectives (pertainyms)
        "adv.all",            // all adverbs
        "noun.Tops",          // unique beginner for nouns
        "noun.act",           // nouns denoting acts or actions
        "noun.animal",        // nouns denoting animals
        "noun.artifact",      // nouns denoting man-made objects
        "noun.attribute",     // nouns denoting attributes of people and objects
        "noun.body",          // nouns denoting body parts
        "noun.cognition",     // nouns denoting cognitive processes and contents
        "noun.communication", // nouns denoting communicative processes and contents
        "noun.event",         // nouns denoting natural events
        "noun.feeling",       // nouns denoting feelings and emotions
        "noun.food",          // nouns denoting foods and drinks
        "noun.group",         // nouns denoting groupings of people or objects
        "noun.location",      // nouns denoting spatial position
        "noun.motive",        // nouns denoting goals
        "noun.object",        // nouns denoting natural objects (not man-made)
        "noun.person",        // nouns denoting people
        "noun.phenomenon",    // nouns denoting natural phenomena
        "noun.plant",         // nouns denoting plants
        "noun.possession",    // nouns denoting possession and transfer of possession
        "noun.process",       // nouns denoting natural processes
        "noun.quantity",      // nouns denoting quantities and units of measure
        "noun.relation",      // nouns denoting relations between people or things or ideas
        "noun.shape",         // nouns denoting two and three dimensional shapes
        "noun.state",         // nouns denoting stable states of affairs
        "noun.substance",     // nouns denoting substances
        "noun.time",          // nouns denoting time and temporal relations
        "verb.body",          // verbs of grooming, dressing and bodily care
        "verb.change",        // verbs of size, temperature change, intensifying, etc.
        "verb.cognition",     // verbs of thinking, judging, analyzing, doubting
        "verb.communication", // verbs of telling, asking, ordering, singing
        "verb.competition",   // verbs of fighting, athletic activities
        "verb.consumption",   // verbs of eating and drinking
        "verb.contact",       // verbs of touching, hitting, tying, digging
        "verb.creation",      // verbs of sewing, baking, painting, performing
        "verb.emotion",       // verbs of feeling
        "verb.motion",        // verbs of walking, flying, swimming
        "verb.perception",    // verbs of seeing, hearing, feeling
        "verb.possession",    // verbs of buying, selling, owning
        "verb.social",        // verbs of political and social activities and events
        "verb.stative",       // verbs of being, having, spatial relations
        "verb.weather",       // verbs of raining, snowing, thawing, thundering
        "adj.ppl",            // participial adjectives
    }

    RELATIONSHIP_POINTER_SYMBOLS = map[string]int {
        // noun relationships
        "!": ANTONYM_RELATIONSHIP,
        "@": HYPERNYM_RELATIONSHIP,
        "@i": INSTANCE_HYPERNYM_RELATIONSHIP,
        "~": HYPONYM_RELATIONSHIP,
        "~i": INSTANCE_HYPONYM_RELATIONSHIP,
        "#m": MEMBER_HOLONYM_RELATIONSHIP,
        "#s": SUBSTANCE_HOLONYM_RELATIONSHIP,
        "#p": PART_HOLONYM_RELATIONSHIP,
        "%m": MEMBER_MERONYM_RELATIONSHIP,
        "%s": SUBSTANCE_MERONYM_RELATIONSHIP,
        "%p": PART_MERONYM_RELATIONSHIP,
        "=": ATTRIBUTE_RELATIONSHIP,
        "+": DERIVATIONALLY_RELATED_FORM_RELATIONSHIP,
        ";": DOMAIN_OF_SYNSET,
        ";c": DOMAIN_OF_SYNSET_TOPIC_RELATIONSHIP,
        "-c": MEMBER_OF_THIS_DOMAIN_TOPIC_RELATIONSHIP,
        ";r": DOMAIN_OF_SYNSET_REGION_RELATIONSHIP,
        "-r": MEMBER_OF_THIS_DOMAIN_REGION_RELATIONSHIP,
        ";u": DOMAIN_OF_SYNSET_USAGE_RELATIONSHIP,
        "-u": MEMBER_OF_THIS_DOMAIN_USAGE_RELATIONSHIP,

        // verb relationships
        // ANTONYM_RELATIONSHIP
        // HYPERNYM_RELATIONSHIP
        // HYPONYM_RELATIONSHIP
        "*": ENTAILMENT_RELATIONSHIP,
        ">": CAUSAL_RELATIONSHIP,
        "^": ALSO_SEE_RELATIONSHIP,
        "$": VERB_GROUP_RELATIONSHIP,
        // DERIVATIONALLY_RELATED_FORM_RELATIONSHIP
        // DOMAIN_OF_SYNSET_TOPIC_RELATIONSHIP,
        // DOMAIN_OF_SYNSET_REGION_RELATIONSHIP,
        // DOMAIN_OF_SYNSET_USAGE_RELATIONSHIP,

        // adjective relationships
        // ANTONYM_RELATIONSHIP
        "&": SIMILAR_TO_RELATIONSHIP,
        "<": PARTICIPLE_OF_VERB_RELATIONSHIP,
        "\\": PERTAINYM_RELATIONSHIP,
        // ATTRIBUTE_RELATIONSHIP
        // ALSO_SEE_RELATIONSHIP
        // DOMAIN_OF_SYNSET_TOPIC_RELATIONSHIP
        // DOMAIN_OF_SYNSET_REGION_RELATIONSHIP
        // DOMAIN_OF_SYNSET_USAGE_RELATIONSHIP

        // adverb relationships
        // ANTONYM_RELATIONSHIP
        // PERTAINYM_RELATIONSHIP
        // DOMAIN_OF_SYNSET_TOPIC_RELATIONSHIP
        // DOMAIN_OF_SYNSET_REGION_RELATIONSHIP
        // DOMAIN_OF_SYNSET_USAGE_RELATIONSHIP
    }
)

// syntactic category / part of speech
const POS_UNSUPPORTED int = 0
const POS_NOUN int = 1
const POS_VERB int = 2
const POS_ADJECTIVE int = 3
const POS_ADVERB int = 4
const POS_ADJECTIVE_SATELLITE int = 5

// relations among synsets
const ANTONYM_RELATIONSHIP int = 10
const HYPERNYM_RELATIONSHIP int = 20
const INSTANCE_HYPERNYM_RELATIONSHIP int = 21
const HYPONYM_RELATIONSHIP int = 30
const INSTANCE_HYPONYM_RELATIONSHIP int = 31
const MEMBER_HOLONYM_RELATIONSHIP int = 40
const SUBSTANCE_HOLONYM_RELATIONSHIP int = 41
const PART_HOLONYM_RELATIONSHIP int = 42
const MEMBER_MERONYM_RELATIONSHIP int = 50
const SUBSTANCE_MERONYM_RELATIONSHIP int = 51
const PART_MERONYM_RELATIONSHIP int = 52
const ATTRIBUTE_RELATIONSHIP int = 60
const DERIVATIONALLY_RELATED_FORM_RELATIONSHIP int = 70
const DOMAIN_OF_SYNSET int = 80
const DOMAIN_OF_SYNSET_TOPIC_RELATIONSHIP int = 90
const MEMBER_OF_THIS_DOMAIN_TOPIC_RELATIONSHIP int = 91
const DOMAIN_OF_SYNSET_REGION_RELATIONSHIP int = 100
const MEMBER_OF_THIS_DOMAIN_REGION_RELATIONSHIP int = 101
const DOMAIN_OF_SYNSET_USAGE_RELATIONSHIP int = 110
const MEMBER_OF_THIS_DOMAIN_USAGE_RELATIONSHIP int = 111
const ENTAILMENT_RELATIONSHIP int = 120
const CAUSAL_RELATIONSHIP int = 130
const ALSO_SEE_RELATIONSHIP int = 140
const VERB_GROUP_RELATIONSHIP int = 150
const SIMILAR_TO_RELATIONSHIP int = 160
const PARTICIPLE_OF_VERB_RELATIONSHIP int = 170
const PERTAINYM_RELATIONSHIP int = 180

// syntactic markers for adjectives
const SYNTACTIC_MARKER_NOT_APPLICABLE int = 0
const SYNTACTIC_MARKER_PREDICATE_POSITION int = 1
const SYNTACTIC_MARKER_PRENOMINAL_POSITION int = 2
const SYNTACTIC_MARKER_IMMEDIATELY_POSTNOMIAL_POSITION int = 3

 func getLemmaKey(lemma string, sense_id int) string {
     return fmt.Sprintf("%s%02d", lemma, sense_id)
 }

func readStoredLemma(s string) string {
    return strings.Replace(s, "_", " ", -1)
}

func writeStoredLemma(s string) string {
    return strings.Replace(s, " ", "_", -1)
}

func oneCharPosTagToPosId(tag string) int {
    switch (tag) {
    case "n":
        return POS_NOUN
    case "v":
        return POS_VERB
    case "a":
        return POS_ADJECTIVE
    case "r":
        return POS_ADVERB
    case "s":
        return POS_ADJECTIVE_SATELLITE
    default:
        return POS_UNSUPPORTED
    }
}
