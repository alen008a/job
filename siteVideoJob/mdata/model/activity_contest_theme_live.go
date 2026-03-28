package model

type ActivityContestThemeLive struct {
	VenueCode string `gorm:"venue_code" json:"venueCode"` // 场馆code
	ContestId string `gorm:"contest_id" json:"contestId"` // 赛事id
}
