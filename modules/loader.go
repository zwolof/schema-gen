package modules

import (
	"encoding/json"
	"fmt"
	"os"

	"go-csitems-parser/models"

	"github.com/baldurstod/vdf"
)

// merge keys at the root level into its own vdf.KeyValue
func MergeKeysAtRootLevel(root *vdf.KeyValue) {
	// key_count := len(kv.GetChilds())
	if root == nil {
		panic("root KeyValue is nil")
	}

	if len(root.GetChilds()) == 0 {
		panic("root KeyValue has no child keys")
	}

	cached := make(map[string][]*vdf.KeyValue)

	for _, item := range root.GetChilds() {
		sectionName := item.Key

		if cached[sectionName] == nil {
			cached[sectionName] = []*vdf.KeyValue{}
		}

		// append the children of the item to the new root value
		cached[sectionName] = append(cached[sectionName], item.GetChilds()...)
	}

	// create a new KeyValue for each section
	newRoot := &vdf.KeyValue{
		Key:   "items_game",
		Value: make([]*vdf.KeyValue, 0),
	}

	for sectionName, items := range cached {
		if len(items) == 0 {
			continue
		}

		// create a new KeyValue with the section name
		newKV := &vdf.KeyValue{
			Key:   sectionName,
			Value: items,
		}
		// add the new KeyValue to the root
		newRoot.Value = append(newRoot.Value.([]*vdf.KeyValue), newKV)
	}

	// replace the root value with the new root value
	root.Value = newRoot.Value
}

func LoadKnifeSkinsMap(path string) map[string][]string {
	fileData, err := os.ReadFile(path)

	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}

	if len(fileData) == 0 {
		panic(fmt.Sprintf("File %s is empty", path))
	}

	result := make(map[string][]string)
	json.Unmarshal(fileData, &result)

	return result
}

func LoadItemsGame(path string) *models.ItemsGame {
	fileData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	vdf := vdf.VDF{}
	parsed := vdf.Parse(fileData)

	kv, _ := parsed.Get("items_game")
	if kv == nil {
		panic("items_game.txt does not contain 'items_game' section")
	}

	MergeKeysAtRootLevel(kv)

	return &models.ItemsGame{
		KeyValue: kv,
	}
}
