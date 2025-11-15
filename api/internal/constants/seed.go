// constants/seeds.go
package constants

const (
	RarityCommon    = "common"
	RarityUncommon  = "uncommon"
	RarityRare      = "rare"
	RarityLegendary = "legendary"
	RarityUnique    = "unique"
)

var ValidRarities = []string{
	RarityCommon,
	RarityUncommon, 
	RarityRare,
	RarityLegendary,
	RarityUnique,
}

// Базовые значения для расчета наград в зависимости от редкости
var RarityMultipliers = map[string]float64{
	RarityCommon:    1.0,
	RarityUncommon:  1.5,
	RarityRare:      2.0,
	RarityLegendary: 3.0,
	RarityUnique:    5.0,
}