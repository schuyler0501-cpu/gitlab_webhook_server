package model

import (
	"time"

	"gorm.io/gorm"
)

// MemberContribution 成员贡献聚合统计表
type MemberContribution struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberEmail  string    `gorm:"type:varchar(255);not null;index" json:"member_email"`
	MemberName   string    `gorm:"type:varchar(255)" json:"member_name"`
	ProjectID    *int      `gorm:"type:integer;index" json:"project_id"`
	ProjectName  string    `gorm:"type:varchar(255)" json:"project_name"`
	StartDate    time.Time `gorm:"type:date;not null;index:idx_period" json:"start_date"`
	EndDate      time.Time `gorm:"type:date;not null;index:idx_period" json:"end_date"`
	CommitCount  int       `gorm:"type:integer;default:0" json:"commit_count"`
	Additions    int       `gorm:"type:integer;default:0" json:"additions"`
	Deletions    int       `gorm:"type:integer;default:0" json:"deletions"`
	NetLines     int       `gorm:"type:integer;generatedAlwaysAs:(additions - deletions) stored" json:"net_lines"`
	TotalChanges int       `gorm:"type:integer;generatedAlwaysAs:(additions + deletions) stored" json:"total_changes"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (MemberContribution) TableName() string {
	return "member_contributions"
}

// MemberLanguageStat 成员语言统计表
type MemberLanguageStat struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberEmail string    `gorm:"type:varchar(255);not null;index" json:"member_email"`
	Language    string    `gorm:"type:varchar(100);not null;index" json:"language"`
	LinesAdded  int       `gorm:"type:integer;default:0" json:"lines_added"`
	LinesRemoved int      `gorm:"type:integer;default:0" json:"lines_removed"`
	FileCount   int       `gorm:"type:integer;default:0" json:"file_count"`
	PeriodStart time.Time `gorm:"type:date;not null;index:idx_period" json:"period_start"`
	PeriodEnd   time.Time `gorm:"type:date;not null;index:idx_period" json:"period_end"`
	ProjectID   *int      `gorm:"type:integer" json:"project_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (MemberLanguageStat) TableName() string {
	return "member_language_stats"
}

// BeforeCreate 创建前钩子 - 确保唯一性
func (m *MemberContribution) BeforeCreate(tx *gorm.DB) error {
	var count int64
	query := tx.Model(&MemberContribution{}).
		Where("member_email = ?", m.MemberEmail).
		Where("start_date = ?", m.StartDate).
		Where("end_date = ?", m.EndDate)
	
	if m.ProjectID != nil {
		query = query.Where("project_id = ?", *m.ProjectID)
	} else {
		query = query.Where("project_id IS NULL")
	}
	
	query.Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// BeforeCreate 创建前钩子 - 确保唯一性
func (m *MemberLanguageStat) BeforeCreate(tx *gorm.DB) error {
	var count int64
	query := tx.Model(&MemberLanguageStat{}).
		Where("member_email = ?", m.MemberEmail).
		Where("language = ?", m.Language).
		Where("period_start = ?", m.PeriodStart).
		Where("period_end = ?", m.PeriodEnd)
	
	if m.ProjectID != nil {
		query = query.Where("project_id = ?", *m.ProjectID)
	} else {
		query = query.Where("project_id IS NULL")
	}
	
	query.Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

