package database

type User struct {
	ID int64 `gorm:"primaryKey"`
}

type Word struct {
	ID   uint   `gorm:"primaryKey"`
	Word string `gorm:"uniqueIndex"`
}

type UserWord struct {
	UserID         int64 `gorm:"primaryKey"`
	WordID         uint  `gorm:"primaryKey"`
	KnowledgeLevel int   `gorm:"check:knowledge_level >= 0 AND knowledge_level <= 5"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
}
