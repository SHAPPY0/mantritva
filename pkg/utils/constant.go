package utils

// const (
// 	ADMIN     UserRole = "ADMIN"
// 	MODERATOR UserRole = "MODERATOR"
// 	USER      UserRole = "USER"
// )

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
	SUGGESTION  FeedType = "SUGGESTION"
	POLL        FeedType = "POLL"
)

const (
	UPVOTE   Vote = "UPVOTE"
	DOWNVOTE Vote = "DOWNVOTE"
)

const (
	PRESIDENT                           Position = "PRESIDENT"
	VICE_PRESIDENT                      Position = "VICE_PRESIDENT"
	PRIME_MINISTER                      Position = "PRIME_MINISTER"
	CABINET_MINISTER                    Position = "CABINET_MINISTER"
	MINISTER_OF_STATE_INDEPENDENT_IC    Position = "MINISTER_OF_STATE_INDEPENDENT_IC"
	MINISTER_OF_STATE                   Position = "MINISTER_OF_STATE"
	ATTORNEY_GENERAL                    Position = "ATTORNEY_GENERAL"
	CABINET_SECRETARY                   Position = "CABINET_SECRETARY"
	SPEAKER                             Position = "SPEAKER"
	DEPUTY_SPEAKER                      Position = "DEPUTY_SPEAKER"
	ASSISTANT_SPEAKER                   Position = "ASSISTANT_SPEAKER"
	CHAIRMAN                            Position = "CHAIRMAN"
	DEPUTY_CHAIRMAN                     Position = "DEPUTY_CHAIRMAN"
	ASSISTANT_CHAIRMAN                  Position = "ASSISTANT_CHAIRMAN"
	CHIEF_MINISTER                      Position = "CHIEF_MINISTER"
	DEPUTY_CHIEF_MINISTER               Position = "DEPUTY_CHIEF_MINISTER"
	MEMBER_OF_PARLIAMENT                Position = "MEMBER_OF_PARLIAMENT"
	MEMBER_OF_PARLIAMENT_RS             Position = "MEMBER_OF_PARLIAMENT_RS"
	GOVERNOR                            Position = "GOVERNOR"
	CHIEF_SECRETARY                     Position = "CHIEF_SECRETARY"
	DIRECTOR_GENERAL                    Position = "DIRECTOR_GENERAL"
	DIRECTOR                            Position = "DIRECTOR"
	ADDITIONAL_SECRETARY                Position = "ADDITIONAL_SECRETARY"
	JOINT_SECRETARY                     Position = "JOINT_SECRETARY"
	UNDER_SECRETARY                     Position = "UNDER_SECRETARY"
	ASSISTANT_SECRETARY                 Position = "ASSISTANT_SECRETARY"
	DEPUTY_SECRETARY                    Position = "DEPUTY_SECRETARY"
	ASSISTANT_DIRECTOR                  Position = "ASSISTANT_DIRECTOR"
	DEPUTY_DIRECTOR                     Position = "DEPUTY_DIRECTOR"
	ASSISTANT_DEPUTY_DIRECTOR           Position = "ASSISTANT_DEPUTY_DIRECTOR"
	ADVOCATE_GENERAL                    Position = "ADVOCATE_GENERAL"
	COMMISSIONER                        Position = "COMMISSIONER"
	DEPUTY_COMMISSIONER                 Position = "DEPUTY_COMMISSIONER"
	ASSISTANT_COMMISSIONER              Position = "ASSISTANT_COMMISSIONER"
	SUB_DIVISIONAL_OFFICER              Position = "SUB_DIVISIONAL_OFFICER"
	DEPUTY_SUB_DIVISIONAL_OFFICER       Position = "DEPUTY_SUB_DIVISIONAL_OFFICER"
	ASSISTANT_SUB_DIVISIONAL_OFFICER    Position = "ASSISTANT_SUB_DIVISIONAL_OFFICER"
	DISTRICT_COLLECTOR                  Position = "DISTRICT_COLLECTOR"
	DEPUTY_DISTRICT_COLLECTOR           Position = "DEPUTY_DISTRICT_COLLECTOR"
	MEMBER_OF_LEGISLATIVE_ASSEMBLY      Position = "MEMBER_OF_LEGISLATIVE_ASSEMBLY"
	MEMBER_OF_LEGISLATIVE_COUNCIL       Position = "MEMBER_OF_LEGISLATIVE_COUNCIL"
	MEMBER_OF_LEGISLATIVE_PARLIAMENT    Position = "MEMBER_OF_LEGISLATIVE_PARLIAMENT"
	MEMBER_OF_LEGISLATIVE_PARLIAMENT_IC Position = "MEMBER_OF_LEGISLATIVE_PARLIAMENT_IC"
	MEMBER_OF_LEGISLATIVE_PARLIAMENT_RS Position = "MEMBER_OF_LEGISLATIVE_PARLIAMENT_RS"
	MAYOR                               Position = "MAYOR"
	DEPUTY_MAYOR                        Position = "DEPUTY_MAYOR"
	MUNICIPAL_COUNCILLORS               Position = "MUNICIPAL_COUNCILLORS"
	STANDING_COMMITTEE_CHAIRPERSONS     Position = "STANDING_COMMITTEE_CHAIRPERSONS"
	ALDERMEN                            Position = "ALDERMEN"
	NAGAR_PALIKA_PRESIDENT              Position = "NAGAR_PALIKA_PRESIDENT"
	NAGAR_PALIKA_VICE_PRESIDENT         Position = "NAGAR_PALIKA_VICE_PRESIDENT"
	WARD_MEMBER                         Position = "WARD_MEMBER"
	MUNICIPAL_COUNCILLOR                Position = "MUNICIPAL_COUNCILLOR"
	NOMINATED_MEMBER                    Position = "NOMINATED_MEMBER"
	LEADER_OF_OPPOSITION                Position = "LEADER_OF_OPPOSITION"
)

var PositionDisplayMap = map[Position]string{
	PRESIDENT:                           "President of India",
	VICE_PRESIDENT:                      "Vice President of India",
	PRIME_MINISTER:                      "Prime Minister",
	CABINET_MINISTER:                    "Cabinet Minister",
	MINISTER_OF_STATE_INDEPENDENT_IC:    "MoS-IC",
	MINISTER_OF_STATE:                   "MoS",
	GOVERNOR:                            "Governor",
	CHIEF_MINISTER:                      "Chief Minister",
	DEPUTY_CHIEF_MINISTER:               "Deputy Chief Minister",
	ATTORNEY_GENERAL:                    "Attorney General",
	CABINET_SECRETARY:                   "Cabinet Secretary",
	SPEAKER:                             "Speaker",
	DEPUTY_SPEAKER:                      "Deputy Speaker",
	ASSISTANT_SPEAKER:                   "Assistant Speaker",
	CHAIRMAN:                            "Chairman",
	DEPUTY_CHAIRMAN:                     "Deputy Chairman",
	ASSISTANT_CHAIRMAN:                  "Assistant Chairman",
	MEMBER_OF_PARLIAMENT:                "MP",
	MEMBER_OF_PARLIAMENT_RS:             "MP-RS",
	CHIEF_SECRETARY:                     "Chief Secretary",
	DIRECTOR_GENERAL:                    "Director General",
	DIRECTOR:                            "Director",
	ADDITIONAL_SECRETARY:                "Additional Secretary",
	JOINT_SECRETARY:                     "Joint Secretary",
	UNDER_SECRETARY:                     "Under Secretary",
	ASSISTANT_SECRETARY:                 "Assistant Secretary",
	DEPUTY_SECRETARY:                    "Deputy Secretary",
	ASSISTANT_DIRECTOR:                  "Assistant Director",
	DEPUTY_DIRECTOR:                     "Deputy Director",
	ASSISTANT_DEPUTY_DIRECTOR:           "Assistant Deputy Director",
	ADVOCATE_GENERAL:                    "Advocate General",
	COMMISSIONER:                        "Commissioner",
	DEPUTY_COMMISSIONER:                 "Deputy Commissioner",
	ASSISTANT_COMMISSIONER:              "Assistant Commissioner",
	SUB_DIVISIONAL_OFFICER:              "Sub-Division Officer",
	DEPUTY_SUB_DIVISIONAL_OFFICER:       "Deputy Sub-Division Officer",
	ASSISTANT_SUB_DIVISIONAL_OFFICER:    "Assistant Sub-Division Officer",
	DISTRICT_COLLECTOR:                  "District Collector",
	DEPUTY_DISTRICT_COLLECTOR:           "Deputy District Collector",
	MEMBER_OF_LEGISLATIVE_ASSEMBLY:      "MLA",
	MEMBER_OF_LEGISLATIVE_COUNCIL:       "MLC",
	MEMBER_OF_LEGISLATIVE_PARLIAMENT:    "MLP",
	MEMBER_OF_LEGISLATIVE_PARLIAMENT_IC: "MLP-IC",
	MEMBER_OF_LEGISLATIVE_PARLIAMENT_RS: "MLP-RS",
	MAYOR:                               "Mayor",
	DEPUTY_MAYOR:                        "Deputy Mayor",
	MUNICIPAL_COUNCILLORS:               "Municipal Councillors",
	STANDING_COMMITTEE_CHAIRPERSONS:     "Standing Committee Chairpersons",
	ALDERMEN:                            "Aldermen",
	NAGAR_PALIKA_PRESIDENT:              "Nagar Palika President",
	NAGAR_PALIKA_VICE_PRESIDENT:         "Nagar Palika Vice-President",
	WARD_MEMBER:                         "Ward Member",
	MUNICIPAL_COUNCILLOR:                "Municipal Councillor",
	NOMINATED_MEMBER:                    "Nominated Member",
	LEADER_OF_OPPOSITION:                "LoP",
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
	KRA_DEVELOPMENTS: "Developments & Welfare Score",
	KRA_GRIEVANCES:   "Grievances Score",
	KRA_CALRECORDS:   "Criminal & Legal Records Score",
	KRA_NEWSPRESSES:  "Public Statements Reception Score",
}

const (
	SUGGESTION_SUBMITTED         SuggestionStatus = "SUBMITTED"
	SUGGESTION_UNDER_REVIEW      SuggestionStatus = "UNDER_REVIEW"
	SUGGESTION_PUBLISHED         SuggestionStatus = "PUBLISHED"
	SUGGESTION_ACKNOWLEDGED      SuggestionStatus = "ACKNOWLEDGED"
	SUGGESTION_UNDER_EVALUATION  SuggestionStatus = "UNDER_EVALUATION"
	SUGGESTION_ADOPTED           SuggestionStatus = "ADOPTED"
	SUGGESTION_PARTIALLY_ADOPTED SuggestionStatus = "PARTIALLY_ADOPTED"
	SUGGESTION_DECLINED          SuggestionStatus = "DECLINED"
	SUGGESTION_REJECTED          SuggestionStatus = "REJECTED"
	SUGGESTION_ARCHIVED          SuggestionStatus = "ARCHIVED"
)

const (
	GRIEVANCE_LOW      GrievanceSeverity = "LOW"
	GRIEVANCE_MEDIUM   GrievanceSeverity = "MEDIUM"
	GRIEVANCE_HIGH     GrievanceSeverity = "HIGH"
	GRIEVANCE_CRITICAL GrievanceSeverity = "CRITICAL"
)

const (
	OTP_PURPOSE_EMAIL_VERIFICATION  OTPPurpose = "EMAIL_VERIFICATION"
	OTP_PURPOSE_MOBILE_VERIFICATION OTPPurpose = "MOBILE_VERIFICATION"
	OTP_PURPOSE_PASSWORD_RESET      OTPPurpose = "PASSWORD_RESET"
)

const (
	OTP_MODE_EMAIL    OTPMode = "EMAIL"
	OTP_MODE_SMS      OTPMode = "SMS"
	OTP_MODE_WHATSAPP OTPMode = "WA"
)

const (
	LB_NAGAR_NIGAM     LocalBodyType = "NAGAR_NIGAM"
	LB_NAGAR_PALIKA    LocalBodyType = "NAGAR_PALIKA"
	LB_NAGAR_PANCHAYAT LocalBodyType = "NAGAR_PANCHAYAT"
	LB_OTHER           LocalBodyType = "OTHER"
)

var LocalBodyTypeDisplayMap = map[LocalBodyType]string{
	LB_NAGAR_NIGAM:     "Nagar Nigam",
	LB_NAGAR_PALIKA:    "Nagar Palika",
	LB_NAGAR_PANCHAYAT: "Nagar Panchayat",
	LB_OTHER:           "Other",
}

const (
	USER              Roles = "USER"
	SUPERADMIN        Roles = "SUPER_ADMIN"
	ADMIN             Roles = "ADMIN"
	STATE_ADMIN       Roles = "STATE_ADMIN"
	LOKSABHA_ADMIN    Roles = "LOKSABHA_ADMIN"
	VIDHANSABHA_ADMIN Roles = "VIDHANSABHA_ADMIN"
	DISTRICT_ADMIN    Roles = "DISTRICT_ADMIN"
	LOCALBODY_ADMIN   Roles = "LOCALBODY_ADMIN"
	WARD_ADMIN        Roles = "WARD_ADMIN"
)

var RolesMap = map[Roles]string{
	USER:              "USER",
	SUPERADMIN:        "SUPER_ADMIN",
	ADMIN:             "ADMIN",
	STATE_ADMIN:       "STATE_ADMIN",
	LOKSABHA_ADMIN:    "LOKSABHA_ADMIN",
	VIDHANSABHA_ADMIN: "VIDHANSABHA_ADMIN",
	DISTRICT_ADMIN:    "DISTRICT_ADMIN",
	LOCALBODY_ADMIN:   "LOCALBODY_ADMIN",
	WARD_ADMIN:        "WARD_ADMIN",
}

const (
	ROLE_SCOPE_SUPERADMIN  RoleScope = "GLOBAL"
	ROLE_SCOPE_COUNTRY     RoleScope = "COUNTRY"
	ROLE_SCOPE_STATE       RoleScope = "STATE"
	ROLE_SCOPE_DISTRICT    RoleScope = "DISTRICT"
	ROLE_SCOPE_LOCALBODY   RoleScope = "LOCALBODY"
	ROLE_SCOPE_WARD        RoleScope = "WARD"
	ROLE_SCOPE_LOKSABHA    RoleScope = "LOKSABHA"
	ROLE_SCOPE_VIDHANSABHA RoleScope = "VIDHANSABHA"
	ROLE_SCOPE_USER        RoleScope = "GLOBAL"
)

var RoleRoleScopeMap = map[Roles]RoleScope{
	USER:              ROLE_SCOPE_USER,
	SUPERADMIN:        ROLE_SCOPE_SUPERADMIN,
	ADMIN:             ROLE_SCOPE_COUNTRY,
	STATE_ADMIN:       ROLE_SCOPE_STATE,
	LOKSABHA_ADMIN:    ROLE_SCOPE_LOKSABHA,
	VIDHANSABHA_ADMIN: ROLE_SCOPE_VIDHANSABHA,
	DISTRICT_ADMIN:    ROLE_SCOPE_DISTRICT,
	LOCALBODY_ADMIN:   ROLE_SCOPE_LOCALBODY,
	WARD_ADMIN:        ROLE_SCOPE_WARD,
}

const (
	RULER_SCOPE_CENTRAL   RulerScopeType = "CENTRAL"
	RULER_SCOPE_STATE     RulerScopeType = "STATE"
	RULER_SCOPE_DISTRICT  RulerScopeType = "DISTRICT"
	RULER_SCOPE_LOCALBODY RulerScopeType = "LOCALBODY"
	RULER_SCOPE_WARD      RulerScopeType = "WARD"
)
