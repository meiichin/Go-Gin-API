package campaign

import "time"

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmout     int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImages struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
}
