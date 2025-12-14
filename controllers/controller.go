package controllers

import (
	"fmt"
	"net/http"
	"options-api/models"
	"options-api/views"
	"strings"
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

// HandleFieldValue - обрабатывает GET запросы для получения значений конкретных полей
func (c *Controller) HandleFieldValue(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}

	// Извлекаем имя поля из URL пути
	path := strings.Trim(r.URL.Path, "/")

	// Пропускаем корневой путь и статические пути
	if path == "" || path == "save" || strings.HasPrefix(path, "static/") {
		return nil
	}

	config, err := c.configManager.LoadConfig()
	if err != nil {
		http.Error(w, "Ошибка чтения конфигурации", http.StatusInternalServerError)
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Ищем поле с таким именем
	for _, field := range config.Fields {
		if field.Name == path {
			// Возвращаем значение поля в текстовом формате
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprint(w, field.Value)
			return nil
		}
	}

	// Если поле не найдено
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, fmt.Sprintf("Поле '%s' не найдено", path))
	return nil
}
