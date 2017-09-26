package models

import (
	"github.com/raggaer/castro/app/database"
)

// TemplateHook Struct used for template hooks
type TemplateHook struct {
	ExtensionId string
	Template    string
}

// GetTemplateHooksByName gets a list of template hooks from database based on the hook name
func GetTemplateHooksByName(name string) ([]TemplateHook, error) {
	// Placeholders for query values
	hooks := []TemplateHook{}

	// Get hooks from database
	rows, err := database.DB.Queryx("SELECT extension_id, template FROM castro_extension_templatehooks WHERE enabled = 1 AND type = ?", name)

	if err != nil {
		return hooks, err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {
		var extensionID, template string
		if err := rows.Scan(&extensionID, &template); err != nil {
			return hooks, err
		}

		hooks = append(hooks, TemplateHook{ExtensionId: extensionID, Template: template})
	}
	return hooks, nil
}
