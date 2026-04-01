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
	PRIMEMINISTER       Position = "Prime Minister"
	PRESIDENT           Position = "President"
	VICEPRESIDENT       Position = "Vice President"
	CABINETMINISTER     Position = "Cabinet Minister"
	CHIEFMINISTER       Position = "Chief Minister"
	DEPUTYCHIEFMINISTER Position = "Deputy Chief Minister"
	OPPOSITIONLEADER    Position = "Opposition Leader"
	MP                  Position = "Member Of Parliament"
	MOS                 Position = "Member Of State"
	MLA                 Position = "Member Of Legislative Assembly"
	RAJYASABHAMEMEBER   Position = "Rajyasabha Member"
)

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

const (
	GRIEVANCE_LOW      GrievanceSeverity = "LOW"
	GRIEVANCE_MEDIUM   GrievanceSeverity = "MEDIUM"
	GRIEVANCE_HIGH     GrievanceSeverity = "HIGH"
	GRIEVANCE_CRITICAL GrievanceSeverity = "CRITICAL"
)
