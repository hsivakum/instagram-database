package models

import (
	"time"
)

type User struct {
	ID               string `gorm:"primaryKey,default:uuid_generate_v4()"`
	Username         string `gorm:"unique"`
	FollowingCount   int64  `gorm:"default:0"`
	FollowersCount   int64  `gorm:"default:0"`
	Bio              string
	PostsCount       int64 `gorm:"default:0"`
	HighlightsCount  int64 `gorm:"default:0"`
	Name             string
	ProfileImageLink string
	IsBusiness       bool
	IsVerified       bool
	Country          string
	Region           string
	CreatedAt        time.Time  `gorm:"autoCreateTime"`
	UpdatedAt        *time.Time `gorm:"autoUpdateTime"`
	DeletedAt        *time.Time `gorm:"index"`
}

type Business struct {
	ID            string `gorm:"primaryKey,default:uuid_generate_v4()"`
	UserID        string `gorm:"index"`
	CityName      string
	Latitude      float64
	Longitude     float64
	StreetAddress string
	ZipCode       int
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `gorm:"index"`
	User          *User      `gorm:"foreignKey:UserID"`
}

type Follower struct {
	FollowerID  string    `gorm:"primaryKey"`
	FollowingID string    `gorm:"primaryKey"`
	FollowedAt  time.Time `gorm:"autoCreateTime"`
	Follower    User      `gorm:"foreignKey:FollowerID"`
	Following   User      `gorm:"foreignKey:FollowingID"`
}

type FollowersActivity struct {
	ID          string    `gorm:"primaryKey"`
	FollowerID  string    `gorm:"index"`
	FollowingID string    `gorm:"index"`
	IsUnfollow  bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Follower    User      `gorm:"foreignKey:FollowerID"`
	Following   User      `gorm:"foreignKey:FollowingID"`
}

type Location struct {
	ID            int64  `gorm:"primaryKey,autoIncrement"`
	HasPublicPage bool   `gorm:"default:false"`
	Name          string `gorm:"unique"`
	Slug          string `gorm:"unique"`
}

type Post struct {
	ID              *int64 `gorm:"primaryKey"`
	UserID          string `gorm:"index"`
	Caption         string
	LikesCount      int64 `gorm:"default:0"`
	CommentsCount   int64 `gorm:"default:0"`
	VideoViewCount  int64 `gorm:"default:0"`
	PrimaryImageURL string
	PrimaryVideoURL string
	LocationID      *int64     `gorm:"index"`
	IsSponsored     bool       `gorm:"default:false"`
	SponsorID       *string    `gorm:"index"`
	URL             string     `gorm:"not null"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"index"`
	User            User       `gorm:"foreignKey:UserID"`
	Location        Location   `gorm:"foreignKey:LocationID"`
	Sponsor         User       `gorm:"foreignKey:SponsorID"`
}

type PostImage struct {
	ID        int64      `gorm:"primaryKey"`
	PostID    int64      `gorm:"index"`
	ImageURL  string     `gorm:"not null"`
	PostOrder int        `gorm:"default:1"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	Post      Post       `gorm:"foreignKey:PostID"`
}

type Highlight struct {
	ID        int64  `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Title     string
	Image     string     `gorm:"not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	User      User       `gorm:"foreignKey:UserID"`
}

type Story struct {
	ID        string     `gorm:"primaryKey"`
	UserID    string     `gorm:"index"`
	MediaURL  string     `gorm:"not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	User      User       `gorm:"foreignKey:UserID"`
}

type StoryView struct {
	StoryID  string    `gorm:"primaryKey"`
	ViewerID string    `gorm:"primaryKey"`
	IsLiked  bool      `gorm:"default:false"`
	ViewedAt time.Time `gorm:"autoCreateTime"`
	Story    Story     `gorm:"foreignKey:StoryID"`
	Viewer   User      `gorm:"foreignKey:ViewerID"`
}

type HighlightsStory struct {
	HighlightID int64     `gorm:"primaryKey"`
	StoryID     string    `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Highlight   Highlight `gorm:"foreignKey:HighlightID"`
	Story       Story     `gorm:"foreignKey:StoryID"`
}

type HighlightsStoryActivity struct {
	HighlightID int64     `gorm:"primaryKey"`
	StoryID     string    `gorm:"primaryKey"`
	IsRemoved   bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Highlight   Highlight `gorm:"foreignKey:HighlightID"`
	Story       Story     `gorm:"foreignKey:StoryID"`
}

type HashTag struct {
	ID        int64      `gorm:"primaryKey"`
	Name      string     `gorm:"not null"`
	CreatedBy *string    `gorm:"index"`
	IsBlocked bool       `gorm:"default:false"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	Creator   User       `gorm:"foreignKey:CreatedBy"`
}

type PostTag struct {
	PostID    int64     `gorm:"primaryKey"`
	TagID     int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Post      Post      `gorm:"foreignKey:PostID"`
	Tag       HashTag   `gorm:"foreignKey:TagID"`
}

type StoryTag struct {
	StoryID   string    `gorm:"primaryKey"`
	TagID     int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Story     Story     `gorm:"foreignKey:StoryID"`
	Tag       HashTag   `gorm:"foreignKey:TagID"`
}

type Block struct {
	UserID    string    `gorm:"primaryKey"`
	BlockedID string    `gorm:"primaryKey"`
	BlockedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID"`
	Blocked   User      `gorm:"foreignKey:BlockedID"`
}

type BlockActivity struct {
	UserID    string     `gorm:"primaryKey"`
	BlockedID string     `gorm:"primaryKey"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	IsBlock   bool       `gorm:"default:true"`
	User      User       `gorm:"foreignKey:UserID"`
	Blocked   User       `gorm:"foreignKey:BlockedID"`
}

type Restrict struct {
	UserID         string     `gorm:"primaryKey"`
	RestrictUserID string     `gorm:"primaryKey"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`
	User           User       `gorm:"foreignKey:UserID"`
	RestrictUser   User       `gorm:"foreignKey:RestrictUserID"`
}

type RestrictActivity struct {
	UserID         string    `gorm:"primaryKey"`
	RestrictUserID string    `gorm:"primaryKey"`
	IsRestrict     bool      `gorm:"default:true"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	User           User      `gorm:"foreignKey:UserID"`
	RestrictUser   User      `gorm:"foreignKey:RestrictUserID"`
}

type Comment struct {
	ID              int64  `gorm:"primaryKey"`
	PostID          int64  `gorm:"index"`
	UserID          string `gorm:"index"`
	ParentCommentID int64
	CommentText     string     `gorm:"not null"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"index"`
	Post            Post       `gorm:"foreignKey:PostID"`
	User            User       `gorm:"foreignKey:UserID"`
}

type CommentLike struct {
	CommentID int64     `gorm:"primaryKey"`
	LikedBy   string    `gorm:"primaryKey"`
	LikedAt   time.Time `gorm:"autoCreateTime"`
	Comment   Comment   `gorm:"foreignKey:CommentID"`
	User      User      `gorm:"foreignKey:LikedBy"`
}

type CommentActivity struct {
	CommentID int64     `gorm:"primaryKey"`
	ActionBy  string    `gorm:"primaryKey"`
	IsLike    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Comment   Comment   `gorm:"foreignKey:CommentID"`
	User      User      `gorm:"foreignKey:ActionBy"`
}
