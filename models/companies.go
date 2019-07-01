package models

type CompaniesJSONItem struct {
	// {"custno":"e24f07e7-5670-40ea-b624-184b5fdc50d9","invoice":53935562,"name":"Amazon Web Services Taiwan Limited_台灣亞馬遜網路服務有限公司","management":null,"product":"為AWS用戶提供技術上的支援與協助 \n運用AWS強大的後端工具與領先的技術，為客戶的需求提供客製化的解決方案 \n依據您的豐富經驗，給AWS研發團隊提供寶貴的建議，協助產品不斷進步。","profile":"關於Amazon Web Services \nAmazon Web Services成立於 2006 年，透過位於美國、澳洲、巴西、中國、德國、愛爾蘭、日本與新加坡的資料中心，提供強大而功能完整的雲端基礎設施平台，其廣泛的服務包含運算、儲存、資料庫、分析、應用與部署服務。全球 190 個國家快速成長的新創公司、大型企業、政府機構等超過 100 萬個客戶目前仰賴 AWS 的服務快速創新、降低 IT 成本以及擴張全球應用。如需了解更多 AWS 相關資訊，請至網站  http:\/\/aws.amazon.com\/tw。","welfare":null}
	Custno     *string
	Invoice    *int64
	Name       *string
	Profile    *string
	Management *string
	Welfare    *string
	Product    *string
}
