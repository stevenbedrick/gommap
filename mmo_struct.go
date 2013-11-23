package metamap

import "time"

/*
Struct that matches the structure of MetaMap's XML output.
*/

type MMOs struct {
	MMO struct {
		CmdLine struct {
			Command	string
			Options struct {
				Count int `xml:"Count,attr"`
				Options	[]Option	`xml:"Option"`
			}
		}
		Negations struct {
			Count int `xml:"Count,attr"`
			Negations []Negation `xml:"Negation"`
		}
		
		Utterances struct {
			Count int `xml:"Count,attr"`
			Utterances []Utterance `xml:"Utterance"`
		}
	}
	ParseTime time.Duration
	RawXML string
}

type Option struct {
	OptName string
	OptValue string
}

type Negation struct {
	NegType string
	NegTrigger string
	NegTriggerPIs struct {
		Count int `xml:"Count,attr"`
		NegTriggerPIs []NegTriggerPI `xml:"NegTriggerPI"`
	}
	
	NegConcepts struct {
		Count int `xml:"Count,attr"`
		NegConcepts []NegConcept `xml:"NegConcept"`
	}
	
	NegConcPIs struct {
		Count int `xml:"Count,attr"`
		NegConcPIs []NegConcPI `xml:"NegConcPI"`
	}
}

type NegTriggerPI struct {
	StartPost int
	Length int
}

type NegConcept struct {
	NegConcCUI string
	NegConcMatched string
}

type NegConcPI struct {
	StartPos int
	Length int
}

type Utterance struct {
	PMID string
	UttSection string
	UttNum int
	UttText string
	UttStartPos int
	UttLength int
	Phrases struct {
		Count int `xml:"Count,attr"`
		Phrases []Phrase `xml:"Phrase"`
	}
}

type Phrase struct {
	PhraseText string
	SyntaxUnits struct {
		Count int `xml:"Count,attr"`
		SyntaxUnits []SyntaxUnit `xml:"SyntaxUnit"`
	}
	PhraseStartPos int
	PhraseLength int
	Candidates struct {
		Total int `xml:"Total,attr"`
		Excluded int `xml:"Excluded,attr"`
		Pruned int `xml:"Pruned,attr"`
		Remaining int `xml:"Remaining"`
		Candidates []Candidate `xml:"Candidate"`
	}
}

type SyntaxUnit struct {
	SyntaxType string
	LexMatch string
	InputMatch string
	LexCat string
	Tokens struct {
		Count int `xml:"Count,attr"`
		Tokens []string `xml:"Token"`
	}
}

type Candidate struct {
	CandidateScore int
	CandidateCUI string
	CandidateMatched string
	CandidatePreferred string
	MatchedWords struct {
		Count int `xml:"Count,attr"`
		MatchedWords []string `xml:"MatchedWord"`
	}
	SemTypes struct {
		Count int `xml:"Count,attr"`
		SemTypes []string `xml:"SemType"`
	}
	MatchMaps struct {
		Count int `xml:"Count,attr"`
		MatchMaps []MatchMap `xml:"MatchMap"`
	}
	IsHead string
	IsOverMatch string
	Sources struct {
		Count int `xml:"Count,attr"`
		Sources []string `xml:"Source"`
	}
	ConceptPIs struct {
		Count int `xml:"Count,attr"`
		ConceptPIs []ConceptPI
	}
	Status int
}

type MatchMap struct {
	TextMatchStart int
	TextMatchEnd int
	ConcMatchStart int
	ConcMatchEnd int
	LexVariation int
}

type ConceptPI struct {
	StartPos int
	Length int
}


