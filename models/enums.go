package models

type CollectibleType int
type StickerEffect int
type StickerType int
type ItemSetType int

const (
	CollectibleTypeUnknown            CollectibleType = -1
	CollectibleTypeServiceMedal       CollectibleType = 0
	CollectibleTypeMapContributor     CollectibleType = 10
	CollectibleTypeMapPin             CollectibleType = 20
	CollectibleTypeOperation          CollectibleType = 30
	CollectibleTypePickEm             CollectibleType = 40
	CollectibleTypeOldPickEm          CollectibleType = 50
	CollectibleTypeFantasyTrophy      CollectibleType = 60
	CollectibleTypeTournamentFinalist CollectibleType = 70
	CollectibleTypePremierSeasonCoin  CollectibleType = 80
	CollectibleTypeYearsOfService     CollectibleType = 90
)

const (
	StickerEffectUnknown    StickerEffect = -1
	StickerEffectNormal     StickerEffect = 0
	StickerEffectHolo       StickerEffect = 1
	StickerEffectFoil       StickerEffect = 2
	StickerEffectGold       StickerEffect = 3
	StickerEffectGlitter    StickerEffect = 4
	StickerEffectLenticular StickerEffect = 5
)

const (
	StickerTypeUnknown   StickerType = -1
	StickerTypeAutograph StickerType = 0
	StickerTypeTeam      StickerType = 1
	StickerTypeEvent     StickerType = 2
)

const (
	ItemSetTypeUnknown   ItemSetType = -1
	ItemSetTypePaintKits ItemSetType = 0
	ItemSetTypeAgents    ItemSetType = 1
	ItemSetTypeStickers  ItemSetType = 2
)
