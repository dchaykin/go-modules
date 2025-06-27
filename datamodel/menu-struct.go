package datamodel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/httpcomm"
)

type Menu struct {
	Menu map[string][]string `json:"menu,omitempty"`
}

type MenuItemConfig struct {
	Name     string           `json:"name"`
	Route    string           `json:"route"`
	Icon     *string          `json:"icon,omitempty"`
	SubItems []MenuItemConfig `json:"items,omitempty"` // only one sub-level is supported by the webclient
}

type MenuConfig struct {
	Items []MenuItemConfig `json:"config"`
	Roles map[string]Menu  `json:"roles"`
}

func (mc *MenuConfig) ReadFromFile(path string, version int) error {
	fullPath := fmt.Sprintf("%s-%03d", path, version)
	jsonData, err := os.ReadFile(fullPath + "/menu-struct.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, mc)
}

func (mc MenuConfig) CreateMenuByRole(userRole string) []MenuItemConfig {
	defMenu, ok := mc.Roles["default"]
	if !ok {
		defMenu = Menu{}
	}
	userRoleMenu, ok := mc.Roles[userRole]
	if !ok {
		userRoleMenu = Menu{}
	}

	menu := mc.mergeMenus(defMenu, userRoleMenu)
	return mc.filterMenuItems(menu)
}

func (mc MenuConfig) mergeMenus(a, b Menu) Menu {
	result := Menu{
		Menu: make(map[string][]string),
	}
	// Duplikat-Tracking je Schlüssel
	seen := make(map[string]map[string]struct{})

	// Schlüssel sammeln (auch leere!) aus beiden Menüs
	allKeys := make(map[string]struct{})
	for k := range a.Menu {
		allKeys[k] = struct{}{}
	}
	for k := range b.Menu {
		allKeys[k] = struct{}{}
	}

	// Für jeden Schlüssel: beide Seiten zusammenführen + Duplikate entfernen
	for key := range allKeys {
		if seen[key] == nil {
			seen[key] = make(map[string]struct{})
		}
		var merged []string
		for _, src := range []map[string][]string{a.Menu, b.Menu} {
			if items, exists := src[key]; exists {
				for _, item := range items {
					if _, exists := seen[key][item]; !exists {
						merged = append(merged, item)
						seen[key][item] = struct{}{}
					}
				}
			}
		}
		// Wichtig: auch leere Listen erhalten
		result.Menu[key] = merged
	}

	return result
}

func (mc MenuConfig) filterMenuItems(menu Menu) []MenuItemConfig {
	var result []MenuItemConfig

	for _, item := range mc.Items {
		// Ist der Haupt-Eintrag erlaubt?
		allowedSubItems, ok := menu.Menu[item.Name]
		if !ok {
			continue
		}

		// Wenn keine SubItems vorhanden → direkt übernehmen
		if len(item.SubItems) == 0 {
			result = append(result, item)
			continue
		}

		// Falls SubItems erlaubt: filtere entsprechend
		allowedSet := make(map[string]struct{}, len(allowedSubItems))
		for _, name := range allowedSubItems {
			allowedSet[name] = struct{}{}
		}

		var filteredSubItems []MenuItemConfig
		for _, sub := range item.SubItems {
			if _, ok := allowedSet[sub.Name]; ok {
				filteredSubItems = append(filteredSubItems, sub)
			}
		}

		// SubItems ggf. ersetzt
		item.SubItems = filteredSubItems
		result = append(result, item)
	}

	return result
}

func GetMenuItemsFromRequest(w http.ResponseWriter, r *http.Request, appName, subPath string) {
	tenant, version, err := httpcomm.GetTenantVersionFromRequest(r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusBadRequest)
		return
	}

	userIdentity, err := auth.GetUserIdentityFromRequest(*r)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusUnauthorized)
		return
	}

	if subPath != "" {
		subPath += "/"
	}

	mc := MenuConfig{}
	err = mc.ReadFromFile("config/"+subPath+tenant, version)
	if err != nil {
		httpcomm.SetResponseError(&w, "", err, http.StatusInternalServerError)
		return
	}

	result := mc.CreateMenuByRole(userIdentity.RoleByApp(appName))

	httpcomm.ServiceResponse{
		Data: result,
	}.WriteData(w, httpcomm.PayloadFormatJSON)
}
