package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

type Field struct {
	Name     string
	Value    interface{}
	Type     string
	Editable bool
}

type PageData struct {
	Fields  []Field
	Success bool
}

var configFile = "options.json"

func main() {
	// Загружаем начальную конфигурацию
	loadConfig()

	// Статические файлы (если потребуются)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Главная страница с формой
	http.HandleFunc("/", homeHandler)

	// Обработчик сохранения
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Читаем текущую конфигурацию
		config, err := readConfig()
		if err != nil {
			http.Error(w, "Ошибка чтения конфигурации", http.StatusInternalServerError)
			return
		}

		// Преобразуем в поля для формы
		var fields []Field
		for k, v := range config {
			fieldType := getFieldType(v)
			fields = append(fields, Field{
				Name:     k,
				Value:    v,
				Type:     fieldType,
				Editable: true,
			})
		}

		// Проверяем параметр успеха
		success := r.URL.Query().Get("success") == "1"

		// Генерируем HTML
		tmpl := template.Must(template.New("index").Parse(htmlTemplate()))
		data := PageData{Fields: fields, Success: success}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Printf("Ошибка генерации HTML: %v\n", err)
			return
		}
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Читаем текущую конфигурацию для получения типов
		currentConfig, err := readConfig()
		if err != nil {
			http.Error(w, "Ошибка чтения конфигурации", http.StatusInternalServerError)
			return
		}

		// Парсим форму
		r.ParseForm()

		// Создаем новую конфигурацию
		newConfig := make(map[string]interface{})

		// Обновляем значения с сохранением типов
		for key, value := range currentConfig {
			formValue := r.Form.Get(key)
			if formValue != "" {
				// Преобразуем строку в соответствующий тип
				convertedValue, err := convertValue(formValue, value)
				if err != nil {
					http.Error(w, fmt.Sprintf("Ошибка преобразования значения для поля %s", key), http.StatusBadRequest)
					return
				}
				newConfig[key] = convertedValue
			} else {
				// Для checkbox (bool), если значение не пришло - ставим false
				if reflect.TypeOf(value).Kind() == reflect.Bool {
					newConfig[key] = false
				} else {
					newConfig[key] = value
				}
			}
		}

		// Проверяем checkbox отдельно, так как они не отправляются, если не отмечены
		for key, value := range currentConfig {
			if reflect.TypeOf(value).Kind() == reflect.Bool {
				// Если checkbox отмечен, значение будет "on"
				if r.Form.Get(key) == "on" {
					newConfig[key] = true
				} else {
					newConfig[key] = false
				}
			}
		}

		// Сохраняем в файл
		err = saveConfig(newConfig)
		if err != nil {
			http.Error(w, "Ошибка сохранения конфигурации", http.StatusInternalServerError)
			return
		}

		// Редирект на главную с индикатором успеха
		http.Redirect(w, r, "/?success=1", http.StatusSeeOther)
	}
}

func readConfig() (map[string]interface{}, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	decoder.UseNumber() // Для корректной обработки чисел

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Преобразуем json.Number в int или float
	for k, v := range config {
		if num, ok := v.(json.Number); ok {
			if i, err := num.Int64(); err == nil {
				config[k] = i
			} else if f, err := num.Float64(); err == nil {
				config[k] = f
			}
		}
	}

	return config, nil
}

func saveConfig(config map[string]interface{}) error {
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(config)
}

func getFieldType(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "text"
	case bool:
		return "checkbox"
	case int, int64, float64:
		return "number"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func convertValue(str string, target interface{}) (interface{}, error) {
	switch target.(type) {
	case string:
		return str, nil
	case bool:
		return str == "on" || str == "true", nil
	case int:
		return strconv.Atoi(str)
	case int64:
		return strconv.ParseInt(str, 10, 64)
	case float64:
		return strconv.ParseFloat(str, 64)
	default:
		// Пробуем определить по содержимому строки
		if val, err := strconv.Atoi(str); err == nil {
			return val, nil
		}
		if val, err := strconv.ParseFloat(str, 64); err == nil {
			return val, nil
		}
		if val, err := strconv.ParseBool(str); err == nil {
			return val, nil
		}
		return str, nil
	}
}

func loadConfig() {
	// Проверяем существование файла
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Создаем файл с дефолтными значениями
		defaultConfig := map[string]interface{}{
			"fieldText": "Текстовое значение",
			"intData":   100500,
			"boolValue": true,
		}

		file, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(defaultConfig)
		fmt.Println("Создан файл конфигурации с настройками по умолчанию")
	}
}

func htmlTemplate() string {
	return `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Настройки</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
            padding: 20px;
        }
        
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
        }
        
        h1 {
            color: #2c3e50;
            margin-bottom: 30px;
            padding-bottom: 15px;
            border-bottom: 2px solid #eee;
        }
        
        .form-group {
            margin-bottom: 25px;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 6px;
            border-left: 4px solid #3498db;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #2c3e50;
            font-size: 16px;
        }
        
        input[type="text"],
        input[type="number"] {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        
        input[type="text"]:focus,
        input[type="number"]:focus {
            outline: none;
            border-color: #3498db;
        }
        
        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        input[type="checkbox"] {
            width: 20px;
            height: 20px;
            cursor: pointer;
        }
        
        .checkbox-label {
            margin-bottom: 0;
            cursor: pointer;
        }
        
        .save-button {
            background: linear-gradient(135deg, #3498db, #2980b9);
            color: white;
            border: none;
            padding: 15px 40px;
            font-size: 18px;
            font-weight: 600;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.3s;
            display: block;
            margin: 40px auto 0;
            width: 200px;
            text-align: center;
        }
        
        .save-button:hover {
            background: linear-gradient(135deg, #2980b9, #3498db);
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(52, 152, 219, 0.3);
        }
        
        .save-button:active {
            transform: translateY(0);
        }
        
        .field-type {
            display: inline-block;
            font-size: 12px;
            color: #7f8c8d;
            background: #ecf0f1;
            padding: 2px 8px;
            border-radius: 10px;
            margin-left: 10px;
        }
        
        .success-message {
            background: #2ecc71;
            color: white;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
            text-align: center;
            animation: fadeIn 0.5s;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        
        .value-display {
            font-family: monospace;
            background: #f1f1f1;
            padding: 8px 12px;
            border-radius: 4px;
            margin-top: 5px;
            font-size: 14px;
            color: #555;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Настройки системы</h1>
        
        {{if .Success}}
        <div class="success-message">
            Настройки успешно сохранены!
        </div>
        {{end}}
        
        <form method="POST" action="/save" id="settingsForm">
            {{range .Fields}}
            <div class="form-group">
                <label for="{{.Name}}">
                    {{.Name}}
                    <span class="field-type">{{.Type}}</span>
                </label>
                
                {{if eq .Type "text"}}
                    <input type="text" 
                           id="{{.Name}}" 
                           name="{{.Name}}" 
                           value="{{.Value}}" 
                           placeholder="Введите текст">
                
                {{else if eq .Type "number"}}
                    <input type="number" 
                           id="{{.Name}}" 
                           name="{{.Name}}" 
                           value="{{.Value}}" 
                           step="any"
                           placeholder="Введите число">
                
                {{else if eq .Type "checkbox"}}
                    <div class="checkbox-group">
                        <input type="checkbox" 
                               id="{{.Name}}" 
                               name="{{.Name}}" 
                               {{if .Value}}checked{{end}}>
                        <label for="{{.Name}}" class="checkbox-label">
                            {{if .Value}}Включено{{else}}Выключено{{end}}
                        </label>
                    </div>
                {{else}}
                    <div class="value-display">
                        {{.Value}}
                    </div>
                {{end}}
            </div>
            {{end}}
            
            <button type="submit" class="save-button">
                💾 Сохранить
            </button>
        </form>
    </div>

    <script>
        // Обновляем текст label для checkbox при изменении
        document.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
            checkbox.addEventListener('change', function() {
                const label = this.nextElementSibling;
                label.textContent = this.checked ? 'Включено' : 'Выключено';
            });
        });
        
        // Подтверждение при перезагрузке страницы с несохраненными изменениями
        let formChanged = false;
        const form = document.getElementById('settingsForm');
        const inputs = form.querySelectorAll('input');
        
        inputs.forEach(input => {
            input.addEventListener('input', () => {
                formChanged = true;
            });
            input.addEventListener('change', () => {
                formChanged = true;
            });
        });
        
        window.addEventListener('beforeunload', (e) => {
            if (formChanged) {
                e.preventDefault();
                e.returnValue = '';
            }
        });
        
        form.addEventListener('submit', () => {
            formChanged = false;
        });
    </script>
</body>
</html>`
}
