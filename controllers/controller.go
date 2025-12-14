package controllers

import (
	"fmt"
	"net/http"
	"options-api/models"
	"options-api/views"
)

type Controller struct {
	configManager *models.ConfigManager
	view          *views.View
}

func NewController(configManager *models.ConfigManager, view *views.View) *Controller {
	return &Controller{
		configManager: configManager,
		view:          view,
	}
}

func (c *Controller) InitializeDefaultConfig() error {
	return c.configManager.CreateDefaultConfigIfNotExists()
}

func (c *Controller) HandleHome(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}

	config, err := c.configManager.LoadConfig()
	if err != nil {
		http.Error(w, "Ошибка чтения конфигурации", http.StatusInternalServerError)
		return fmt.Errorf("failed to load config: %w", err)
	}

	fields := c.configManager.GetFields(config)
	success := r.URL.Query().Get("success") == "1"

	if err := c.view.RenderHome(w, fields, success); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}

func (c *Controller) HandleSave(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}

	currentConfig, err := c.configManager.LoadConfig()
	if err != nil {
		http.Error(w, "Ошибка чтения конфигурации", http.StatusInternalServerError)
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return fmt.Errorf("failed to parse form: %w", err)
	}

	updatedConfig, err := c.configManager.UpdateConfigFromForm(currentConfig, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("failed to update config: %w", err)
	}

	if err := c.configManager.SaveConfig(updatedConfig); err != nil {
		http.Error(w, "Ошибка сохранения конфигурации", http.StatusInternalServerError)
		return fmt.Errorf("failed to save config: %w", err)
	}

	http.Redirect(w, r, "/?success=1", http.StatusSeeOther)
	return nil
}
