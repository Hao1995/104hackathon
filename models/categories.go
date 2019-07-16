package models

type CategoriesFields struct {
	Name string
	File string
}

type CategoriesJSONItem struct {
	// 1001001000,軟體及網路相關業,從事提供顧客特定需要所設計之軟硬體搭配之系統或軟體技術提供系統分析及設計或網路相關之行業。,否
	ID   *int64
	Name *string
	Desc *string
	Hide *string
}
