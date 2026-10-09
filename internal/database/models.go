package database

import "time"

// BanStatus represents a user's ban status.
type BanStatus struct {
	IsBanned    bool      `bson:"is_banned"`
	BanDuration int       `bson:"ban_duration"`
	BannedOn    time.Time `bson:"banned_on"`
	BanReason   string    `bson:"ban_reason"`
}

// User represents a Telegram user.
type User struct {
	ID        int64     `bson:"_id"`
	JoinDate  time.Time `bson:"join_date"`
	Notif     bool      `bson:"notif"`
	BanStatus BanStatus `bson:"ban_status"`
}

// TicketMessage represents a message within a ticket.
type TicketMessage struct {
	FromUserID int64     `bson:"from_user_id"`
	Text       string    `bson:"text"`
	MessageID  int       `bson:"message_id"`
	Timestamp  time.Time `bson:"timestamp"`
	IsAdmin    bool      `bson:"is_admin"`
}

// Ticket represents a support or feedback ticket.
type Ticket struct {
	ID        string          `bson:"_id"`
	UserID    int64           `bson:"user_id"`
	Status    string          `bson:"status"` // open/closed/resolved
	Messages  []TicketMessage `bson:"messages"`
	CreatedAt time.Time       `bson:"created_at"`
	ClosedAt  *time.Time      `bson:"closed_at,omitempty"`
	Rating    int             `bson:"rating"`
	Labels    []string        `bson:"labels"`
}

// Clone represents a cloned bot instance.
type Clone struct {
	BotToken    string    `bson:"bot_token"`
	OwnerID     int64     `bson:"owner_id"`
	BotUsername string    `bson:"bot_username"`
	CreatedAt   time.Time `bson:"created_at"`
	IsActive    bool      `bson:"is_active"`
}

// Stats represents overall bot statistics.
type Stats struct {
	TotalUsers      int64   `bson:"total_users"`
	TotalMessages   int64   `bson:"total_messages"`
	ActiveToday     int64   `bson:"active_today"`
	AvgResponseTime float64 `bson:"avg_response_time"`
}
