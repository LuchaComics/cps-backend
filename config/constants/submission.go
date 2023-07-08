package constants

var SubmissionPublisherNames = map[int8]string{
	1: "Other (Please specify)",
	2: "DC",
	3: "Marvel",
	4: "Image",
	5: "Studiocomix",
	6: "Lucha",
	7: "Boom! Studios",
	8: "Dark Horse Comics",
	9: "IDW",
}

const (
	SubmissionPublisherNameOther = 1
)

var SubmissionOverallLetterGrades = map[string]string{
	"pr": "Poor",
	"fr": "Fair",
	"gd": "Good",
	"vg": "Very good",
	"fn": "Fine",
	"vf": "Very Fine",
	"nm": "Near Mint",
	"PR": "Poor",
	"FR": "Fair",
	"GD": "Good",
	"VG": "Very good",
	"FN": "Fine",
	"VF": "Very Fine",
	"NM": "Near Mint",
}

var SubmissionSpecialDetails = map[int8]string{
	1: "Other",
	2: "Regular Edition",
	3: "Direct Edition",
	4: "Newsstand Edition",
	5: "Variant Cover",
	6: "Canadian Price Variant",
	7: "Facsimile",
	8: "Reprint",
}
