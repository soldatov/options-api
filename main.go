package main

import (
	"fmt"
	"log"
	"net/http"
	"options-api/controllers"
	"options-api/models"
	"options-api/views"
)

var configFile = "options.json"

func main() {
	// Инициализация MVC компонентов
	configManager := models.NewConfigManager(configFile)

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

	// Главная страница с формой
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := controller.HandleHome(w, r); err != nil {
			log.Printf("Error handling home request: %v", err)
		}
	})

	// Обработчик сохранения
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if err := controller.HandleSave(w, r); err != nil {
			log.Printf("Error handling save request: %v", err)
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
