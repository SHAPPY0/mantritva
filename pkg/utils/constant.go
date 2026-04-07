package utils

const (
	ADMIN     UserRole = "ADMIN"
	MODERATOR UserRole = "MODERATOR"
	USER      UserRole = "USER"
)

const (
	PolicyAnnouncement StatementType = "PolicyAnnouncement"
	Reaction           StatementType = "Reaction"
	Criticism          StatementType = "Criticism"
	Defense            StatementType = "Defense"
	ElectionPromise    StatementType = "ElectionPromise"
	OtherStatement     StatementType = "Other"
)

const (
	Economy          Subject = "Economy"
	ForeignPolicy    Subject = "ForeignPolicy"
	SocialIssue      Subject = "SocialIssue"
	ElectionCampaign Subject = "ElectionCampaign"
	OppositionAttack Subject = "OppositionAttack"
	OtherSubject     Subject = "Other"
)

const (
	Supportive StatementTone = "Supportive"
	Critical   StatementTone = "Critical"
	Aggressive StatementTone = "Aggressive"
	Diplomatic StatementTone = "Diplomatic"
	Emotional  StatementTone = "Emotional"
	Other      StatementTone = "Other"
)

const (
	Draft     PostStatus = "Draft"
	Pending   PostStatus = "Pending"
	Published PostStatus = "Published"
	Archived  PostStatus = "Archived"
)

const (
	NEWSPRESS   FeedType = "NEWSPRESS"
	DEVELOPMENT FeedType = "DEVELOPMENT & WELFARE"
	GRIEVANCE   FeedType = "GRIEVANCE"
)

const (
	UPVOTE   Vote = "UPVOTE"
	DOWNVOTE Vote = "DOWNVOTE"
)

const (
	PRESIDENT                     Position = "PRESIDENT"
	VICE_PRESIDENT                Position = "VICE_PRESIDENT"
	PRIME_MINISTER                Position = "PRIME_MINISTER"
	CABINET_MINISTER              Position = "CABINET_MINISTER"
	MINISTER_OF_STATE_INDEPENDENT Position = "MINISTER_OF_STATE_INDEPENDENT"
	MINISTER_OF_STATE             Position = "MINISTER_OF_STATE"
	MP_LS                         Position = "MP_LS"
	MP_RS                         Position = "MP_RS"
	GOVERNOR                      Position = "GOVERNOR"
	CHIEF_MINISTER                Position = "CHIEF_MINISTER"
	DEPUTY_CHIEF_MINISTER         Position = "DEPUTY_CHIEF_MINISTER"
	S_CABINET_MINISTER            Position = "STATE_CABINET_MINISTER"
	S_MINISTER_OF_STATE           Position = "STATE_MINISTER_OF_STATE"
	MLA                           Position = "MLA"
	MLC                           Position = "MLC"
)

var PositionDisplayMap = map[Position]string{
	PRESIDENT:                     "President of India",
	VICE_PRESIDENT:                "Vice President of India",
	PRIME_MINISTER:                "Prime Minister",
	CABINET_MINISTER:              "Cabinet Minister",
	MINISTER_OF_STATE_INDEPENDENT: "MoS-IC",
	MINISTER_OF_STATE:             "MoS",
	MP_LS:                         "MP",
	MP_RS:                         "MP-RS",
	GOVERNOR:                      "Governor",
	CHIEF_MINISTER:                "Chief Minister",
	DEPUTY_CHIEF_MINISTER:         "Deputy Chief Minister",
	S_CABINET_MINISTER:            "State Cabinet Minister",
	S_MINISTER_OF_STATE:           "State MoS",
	MLA:                           "MLA",
	MLC:                           "MLC",
}

const (
	UNDERREVIEW GrievanceStatus = "UNDER_REVIEW"
	OPEN        GrievanceStatus = "OPEN"
	REJECTED    GrievanceStatus = "REJECTED"
	INPROGRESS  GrievanceStatus = "IN_PROGRESS"
	INVALID     GrievanceStatus = "INVALID"
	CLOSED      GrievanceStatus = "CLOSED"
	RESOLVED    GrievanceStatus = "RESOLVED"
	STAYED      GrievanceStatus = "STAYED"
	NOUPDATE    GrievanceStatus = "NO_UPDATE"
)

const (
	ANNOUNCED      DevelopmentStatus = "ANNOUNCED"
	APPROVED       DevelopmentStatus = "APPROVED"
	ENROLLMENTOPEN DevelopmentStatus = "ENROLLMENTOPEN"
	ONGOING        DevelopmentStatus = "ONGOING"
	ONHOLD         DevelopmentStatus = "ONHOLD"
	POSTPONED      DevelopmentStatus = "POSTPONED"
	// COMPLETED DevelopmentStatus = "COMPLETED"
	CANCELLED DevelopmentStatus = "CANCELLED"
)

const (
	KRA_DEVELOPMENTS KRA = "developments"
	KRA_GRIEVANCES   KRA = "grievances"
	KRA_CALRECORDS   KRA = "cal_records"
	KRA_NEWSPRESSES  KRA = "news_and_presses"
)

var KraDisplayMap = map[KRA]string{
	KRA_DEVELOPMENTS: "Developments & Welfare",
	KRA_GRIEVANCES:   "Grievances",
	KRA_CALRECORDS:   "Criminal & Legal Records",
	KRA_NEWSPRESSES:  "News & Presses",
}

const (
	GRIEVANCE_LOW      GrievanceSeverity = "LOW"
	GRIEVANCE_MEDIUM   GrievanceSeverity = "MEDIUM"
	GRIEVANCE_HIGH     GrievanceSeverity = "HIGH"
	GRIEVANCE_CRITICAL GrievanceSeverity = "CRITICAL"
)
