package mini

import (
	"metamap"
	"time"
)

/*
Struct defining the "minimal useful" subset of the full MetaMap output, with annotations for useful JSON output.
*/

	
type MetaMapping struct {
	Phrases []Phrase `json:"phrases"`
	ParseTime time.Duration `json:"parse_time"`
}

type Phrase struct {
	Text string `json:"text"`
	Offsets struct {
		Start int `json:"offset_start"`
		Length int `json:"offset_length"`
	} `json:"offsets"`
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	CUI string `json:"cui"`
	ConceptName string `json:"concept_name"`
	PreferredName string `json:"preferred_name"`
	SemanticTypes []string `json:"semantic_types"`
}

func FromFullMMO(full *metamap.MMOs) *MetaMapping {
	to_ret := &MetaMapping{}
	for _, utt := range(full.MMO.Utterances.Utterances) {
		// for each utterance:
		for _, phr := range(utt.Phrases.Phrases) {
			if phr.Candidates.Total > 0 {
				// for each phrase with mappings:
				// set up a new phrase
				this_p := &Phrase{}
				this_p.Text = phr.PhraseText
				this_p.Offsets.Start = phr.PhraseStartPos
				this_p.Offsets.Length = phr.PhraseLength
				for _, cand := range(phr.Candidates.Candidates) {
					this_m := &Mapping{}
					this_m.CUI = cand.CandidateCUI
					this_m.ConceptName = cand.CandidateMatched
					this_m.PreferredName = cand.CandidatePreferred
					for _, st := range(cand.SemTypes.SemTypes) {
						full_st := SemanticTypeMap()[st]
						// full_st := st
						this_m.SemanticTypes = append(this_m.SemanticTypes, full_st)
					}
					// this_m.SemanticTypes = cand.SemTypes.SemTypes
					this_p.Mappings = append(this_p.Mappings, *this_m)
				}
				to_ret.Phrases = append(to_ret.Phrases, *this_p)
			}
		}
	}
	to_ret.ParseTime = full.ParseTime
	return to_ret
}

/*
Table from srfdef.txt
*/

func SemanticTypeMap() map[string]string {
	return map[string]string{
		"ortf": "Organ or Tissue Function",
		"PT": "part_of",
		"AD": "adjacent_to",
		"RO": "result_of",
		"food": "Food",
		"chvs": "Chemical Viewed Structurally",
		"hcpp": "Human-caused Phenomenon or Process",
		"bodm": "Biomedical or Dental Material",
		"EX": "exhibits",
		"qlco": "Qualitative Concept",
		"diap": "Diagnostic Procedure",
		"menp": "Mental Process",
		"enty": "Entity",
		"eehu": "Environmental Effect of Humans",
		"BR": "branch_of",
		"CM": "complicates",
		"aggp": "Age Group",
		"CS": "consists_of",
		"opco": "Organophosphorus Compound",
		"FR": "functionally_related_to",
		"blor": "Body Location or Region",
		"II": "issue_in",
		"CW": "co-occurs_with",
		"mamm": "Mammal",
		"IG": "ingredient_of",
		"LO": "location_of",
		"qnco": "Quantitative Concept",
		"phsu": "Pharmacologic Substance",
		"celf": "Cell Function",
		"topp": "Therapeutic or Preventive Procedure",
		"inbe": "Individual Behavior",
		"SP": "spatially_related_to",
		"sbst": "Substance",
		"clas": "Classification",
		"fngs": "Fungus",
		"geoa": "Geographic Area",
		"lbtr": "Laboratory or Test Result",
		"PC": "precedes",
		"lang": "Language",
		"ffas": "Fully Formed Anatomical Structure",
		"bdsy": "Body System",
		"comd": "Cell or Molecular Dysfunction",
		"moft": "Molecular Function",
		"elii": "Element, Ion, or Isotope",
		"emod": "Experimental Model of Disease",
		"resa": "Research Activity",
		"TB": "tributary_of",
		"CO": "carries_out",
		"mbrt": "Molecular Biology Research Technique",
		"PR": "physically_related_to",
		"inpr": "Intellectual Product",
		"dsyn": "Disease or Syndrome",
		"DE": "degree_of",
		"MF": "manifestation_of",
		"grpa": "Group Attribute",
		"IC": "interconnects",
		"cell": "Cell",
		"chem": "Chemical",
		"nsba": "Neuroreactive Substance or Biogenic Amine",
		"horm": "Hormone",
		"phob": "Physical Object",
		"carb": "Carbohydrate",
		"lipd": "Lipid",
		"cgab": "Congenital Abnormality",
		"PS": "produces",
		"humn": "Human",
		"imft": "Immunologic Factor",
		"anim": "Animal",
		"acty": "Activity",
		"antb": "Antibiotic",
		"TV": "traverses",
		"CA": "causes",
		"podg": "Patient or Disabled Group",
		"orgm": "Organism",
		"hops": "Hazardous or Poisonous Substance",
		"anab": "Anatomical Abnormality",
		"DS": "disrupts",
		"drdd": "Drug Delivery Device",
		"amph": "Amphibian",
		"enzy": "Enzyme",
		"sosy": "Sign or Symptom",
		"anst": "Anatomical Structure",
		"nusq": "Nucleotide Sequence",
		"orch": "Organic Chemical",
		"dora": "Daily or Recreational Activity",
		"rnlw": "Regulation or Law",
		"celc": "Cell Component",
		"edac": "Educational Activity",
		"ME": "measurement_of",
		"ftcn": "Functional Concept",
		"rcpt": "Receptor",
		"evnt": "Event",
		"shro": "Self-help or Relief Organization",
		"pros": "Professional Society",
		"EV": "evaluation_of",
		"bacs": "Biologically Active Substance",
		"MS": "measures",
		"gngm": "Gene or Genome",
		"MT": "method_of",
		"US": "uses",
		"TR": "temporally_related_to",
		"IN": "indicates",
		"clnd": "Clinical Drug",
		"hlca": "Health Care Activity",
		"mnob": "Manufactured Object",
		"crbs": "Carbohydrate Sequence",
		"cnce": "Conceptual Entity",
		"orga": "Organism Attribute",
		"SR": "surrounds",
		"nnon": "Nucleic Acid, Nucleoside, or Nucleotide",
		"bpoc": "Body Part, Organ, or Organ Component",
		"AN": "analyzes",
		"IS": "isa",
		"tisu": "Tissue",
		"lbpr": "Laboratory Procedure",
		"euka": "Eukaryote",
		"bsoj": "Body Space or Junction",
		"inch": "Inorganic Chemical",
		"phsf": "Physiologic Function",
		"medd": "Medical Device",
		"gora": "Governmental or Regulatory Activity",
		"DF": "developmental_form_of",
		"aapp": "Amino Acid, Peptide, or Protein",
		"virs": "Virus",
		"OC": "occurs_in",
		"ocac": "Occupational Activity",
		"CR": "conceptually_related_to",
		"DI": "diagnoses",
		"tmco": "Temporal Concept",
		"plnt": "Plant",
		"PO": "process_of",
		"bmod": "Biomedical Occupation or Discipline",
		"DO": "derivative_of",
		"neop": "Neoplastic Process",
		"rept": "Reptile",
		"clna": "Clinical Attribute",
		"MN": "manages",
		"genf": "Genetic Function",
		"socb": "Social Behavior",
		"CT": "contains",
		"CN": "connected_to",
		"fndg": "Finding",
		"amas": "Amino Acid Sequence",
		"eico": "Eicosanoid",
		"patf": "Pathologic Function",
		"IW": "interacts_with",
		"mcha": "Machine Activity",
		"bact": "Bacterium",
		"mosq": "Molecular Sequence",
		"PA": "practices",
		"AW": "associated_with",
		"TS": "treats",
		"bdsu": "Body Substance",
		"resd": "Research Device",
		"bird": "Bird",
		"emst": "Embryonic Structure",
		"chvf": "Chemical Viewed Functionally",
		"idcn": "Idea or Concept",
		"vtbt": "Vertebrate",
		"acab": "Acquired Abnormality",
		"ocdi": "Occupation or Discipline",
		"popg": "Population Group",
		"vita": "Vitamin",
		"AF": "affects",
		"arch": "Archaeon",
		"grup": "Group",
		"bhvr": "Behavior",
		"PV": "prevents",
		"PP": "property_of",
		"famg": "Family Group",
		"PE": "performs",
		"prog": "Professional or Occupational Group",
		"spco": "Spatial Concept",
		"inpo": "Injury or Poisoning",
		"npop": "Natural Phenomenon or Process",
		"AE": "assesses_effect_of",
		"orgf": "Organism Function",
		"biof": "Biologic Function",
		"irda": "Indicator, Reagent, or Diagnostic Aid",
		"phpr": "Phenomenon or Process",
		"BA": "brings_about",
		"fish": "Fish",
		"orgt": "Organization",
		"mobd": "Mental or Behavioral Dysfunction",
		"strd": "Steroid",
		"hcro": "Health Care Related Organization",
		"CP": "conceptual_part_of",
	}
}