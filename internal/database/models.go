package database

type Word struct {
	ID   uint   `gorm:"primaryKey"`
	Word string `gorm:"uniqueIndex"`
}

type UserWord struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint
	WordID         uint
	KnowledgeLevel int
	Word           Word `gorm:"foreignKey:WordID"`
}