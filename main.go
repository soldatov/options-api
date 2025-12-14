package main

import (
	"fmt"
	"log"
	"net/http"
	"options-api/controllers"
	"options-api/models"
	"options-api/views"
	"os"
)

func getConfigFile() string {
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		return configFile
	}
	return "options.json"
}

func main() {
	// Инициализация MVC компонентов
	configManager := models.NewConfigManager(getConfigFile())

	view, err := views.NewView()
	if err != nil {
		log.Fatalf("Failed to create view: %v", err)
	}

	controller := controllers.NewController(configManager, view)

	// Загружаем начальную конфигурацию
	if err := controller.InitializeDefaultConfig(); err != nil {
		log.Printf("Warning: Failed to initialize default config: %v", err)
	}

	// Статические файлы (если потребуются)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Универсальный обработчик для всех запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			// Главная страница
			if err := controller.HandleHome(w, r); err != nil {
				log.Printf("Error handling home request: %v", err)
			}
		case "/save":
			// Обработчик сохранения
			if err := controller.HandleSave(w, r); err != nil {
				log.Printf("Error handling save request: %v", err)
			}
		default:
			// Динамические эндпоинты для полей
			if err := controller.HandleFieldValue(w, r); err != nil {
				log.Printf("Error handling field value request: %v", err)
			}
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
