package modules

import "go-csitems-parser/models"

func MapCustomStickers(custom_stickers *[]models.CustomStickers) map[string]models.SchemaCustomSticker {
	rarity_map := make(map[string]models.SchemaCustomSticker, 0)

	for _, rarity := range *custom_stickers {
		// Add to the map
		rarity_map[rarity.GeneratedId] = models.SchemaCustomSticker{
			Group: rarity.Group,
			Name:  rarity.Name,
			Count: rarity.Count,
		}
	}

	return rarity_map
}

func MapCollectibles(collectibles *[]models.Collectible) map[int]models.SchemaGenericeMap {
	collectible_map := make(map[int]models.SchemaGenericeMap)

	for _, collectible := range *collectibles {

		// Add to the map
		collectible_map[collectible.DefinitionIndex] = models.SchemaGenericeMap{
			MarketHashName: collectible.MarketHashName,
			Rarity:         collectible.Rarity,
			Image:          collectible.ImageInventory,
		}
	}

	return collectible_map
}

func MapKeychains(keychains *[]models.Keychain) map[int]SchemaItemWithImage {
	keychain_map := make(map[int]SchemaItemWithImage)

	for _, keychain := range *keychains {
		// Add to the map
		keychain_map[keychain.DefinitionIndex] = SchemaItemWithImage{
			MarketHashName: keychain.MarketHashName,
			Image:          keychain.ImageInventory,
		}
	}

	return keychain_map
}

func MapAgents(agents *[]models.PlayerAgent) map[int]models.SchemaGenericeMap {
	agent_map := make(map[int]models.SchemaGenericeMap)

	for _, agent := range *agents {
		// Add to the map
		agent_map[agent.DefinitionIndex] = models.SchemaGenericeMap{
			MarketHashName: agent.MarketHashName,
			Rarity:         agent.Rarity,
			Image:          agent.ImageInventory,
		}
	}

	return agent_map
}

func MapMusicKits(music_kits *[]models.MusicKit) map[int]models.SchemaGenericeMap {
	music_kit_map := make(map[int]models.SchemaGenericeMap)

	for _, music_kit := range *music_kits {
		// Add to the map
		music_kit_map[music_kit.DefinitionIndex] = models.SchemaGenericeMap{
			MarketHashName: music_kit.MarketHashName,
			Rarity:         "rare", // Music kits are always rare
			Image:          music_kit.ImageInventory,
		}
	}

	return music_kit_map
}

func MapContainers(
	weapon_cases *[]models.WeaponCase,
	souvenir_packages *[]models.SouvenirPackage,
	sticker_capsules *[]models.StickerCapsule,
	pin_capsules *[]models.StickerCapsule,
) map[int]string {
	container_map := make(map[int]string)
	// Add weapon cases to the map
	for _, weapon_case := range *weapon_cases {
		container_map[weapon_case.DefinitionIndex] = weapon_case.MarketHashName
	}

	// Add souvenir packages to the map
	for _, souvenir_package := range *souvenir_packages {
		container_map[souvenir_package.DefinitionIndex] = souvenir_package.MarketHashName
	}

	// Add sticker capsules to the map
	for _, sticker_capsule := range *sticker_capsules {
		container_map[sticker_capsule.DefinitionIndex] = sticker_capsule.MarketHashName
	}

	// Add pin capsules to the map
	for _, pin_capsule := range *pin_capsules {
		container_map[pin_capsule.DefinitionIndex] = pin_capsule.MarketHashName
	}

	return container_map
}

type SchemaItemWithImage struct {
	MarketHashName string `json:"market_hash_name"`
	Image          string `json:"image"`
}

func MapStickerKits(sticker_kits *[]models.StickerKit) map[int]SchemaItemWithImage {
	sticker_kit_map := make(map[int]SchemaItemWithImage)

	for _, sticker_kit := range *sticker_kits {
		// Add to the map
		sticker_kit_map[sticker_kit.DefinitionIndex] = SchemaItemWithImage{
			MarketHashName: sticker_kit.MarketHashName,
			Image:          sticker_kit.StickerMaterial,
		}
	}

	return sticker_kit_map
}
