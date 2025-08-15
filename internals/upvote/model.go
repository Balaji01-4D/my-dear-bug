package upvote

import "time"

// Upvote represents a single user's upvote for a confession.
//
// Uniqueness constraints:
// - (confession_id, ip_hash) unique to prevent multiple votes per IP
// - (confession_id, client_hash) unique to prevent multiple votes per device cookie
// These are best-effort for an anonymous site and help mitigate dynamic IP changes and abuse.
type Upvote struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Participate in two composite unique indexes for IP and client hashing
	ConfessionID uint      `gorm:"uniqueIndex:idx_upvote_conf_ip;uniqueIndex:idx_upvote_conf_client" json:"confessionId"`
	IPHash       string    `gorm:"size:64;uniqueIndex:idx_upvote_conf_ip" json:"-"`
	ClientHash   string    `gorm:"size:64;uniqueIndex:idx_upvote_conf_client" json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
