package parsers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseCollections(
	ctx context.Context,
	ig *models.ItemsGame,
	sv []models.SouvenirPackage,
	cs []models.WeaponCase,
	t *modules.Translator,
) []models.Collection {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	item_sets, err := ig.Get("item_sets")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get item_sets from items_game.txt")
		return nil
	}

	var collections []models.Collection
	for _, s := range item_sets.GetChilds() {
		name, _ := s.GetString("name")

		if strings.Contains(name, "_characters") {
			// Skip if the name contains "_characters"
			continue
		}

		current := models.Collection{
			Key:  s.Key,
			Name: modules.GenerateMarketHashName(t, name, nil, "collection"),
			Image: fmt.Sprintf("econ/set_icons/%s", s.Key),
		}

		// Check if any weapon case matches this item set
		for _, wpncase := range cs {
			if wpncase.ItemSetId == nil || *wpncase.ItemSetId != current.Key {
				continue
			}

			current.HasCrate = true
			break
		}

		// Check if any souvenir package matches this item set
		for _, sv_pkg := range sv {
			if sv_pkg.ItemSetId == nil || *sv_pkg.ItemSetId != current.Key {
				continue
			}
			current.HasSouvenir = true
			break
		}

		// We're done here, add the current item set to the list
		collections = append(collections, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' collections in %s", len(collections), duration)

	return collections
}
