package models

type Criterion struct {
	BaseModel
	Name string `json:"name"`

	Vectors []Vector `json:"vectors" gorm:"foreignKey:CriterionID"`
}
