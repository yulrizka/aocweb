package aocweb

import (
	"encoding/gob"
	"time"

	"gorm.io/gorm"
)

type submissionType string

const (
	unspecified submissionType = "unspecified"
	githubFile  submissionType = "github-file"
	paste       submissionType = "paste"
)

type Submission struct {
	gorm.Model
	Point       int
	Title       string
	Language    string
	Type        submissionType
	Solution1   string `gorm:"size:256"`
	Solution2   string `gorm:"size:256"`
	Description string `gorm:"size:256"`
	User        User
	UserID      int
}

type User struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size: 100"`
	Github    string         `gorm:"size: 100;index:user_github_uniq,unique"`
}

func init() {
	gob.Register(User{})
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}, &Submission{}); err != nil {
		panic(err)
	}
}
