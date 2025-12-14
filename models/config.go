package models

import (
	"encoding/json"
	"fmt"
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

type ConfigField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Config struct {
	Fields []ConfigField `json:"fields"`
}

type ConfigManager struct {
	configFile string
}

func NewConfigManager(configFile string) *ConfigManager {
	return &ConfigManager{configFile: configFile}
}

func (cm *ConfigManager) LoadConfig() (*Config, error) {
	config, err := cm.readConfigFromFile()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (cm *ConfigManager) SaveConfig(config *Config) error {
	return cm.saveConfigToFile(config)
}

func (cm *ConfigManager) CreateDefaultConfigIfNotExists() error {
	if _, err := os.Stat(cm.configFile); os.IsNotExist(err) {
		defaultConfig := &Config{
			Fields: []ConfigField{
				{Name: "fieldText", Value: "Текстовое значение"},
				{Name: "intData", Value: 100500},
				{Name: "boolValue", Value: true},
			},
		}

		file, err := os.Create(cm.configFile)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultConfig); err != nil {
			return err
		}
		fmt.Println("Создан файл конфигурации с настройками по умолчанию")
	}
	return nil
}

func (cm *ConfigManager) GetFields(config *Config) []Field {
	var fields []Field
	for _, field := range config.Fields {
		fieldType := getFieldType(field.Value)
		fields = append(fields, Field{
			Name:     field.Name,
			Value:    field.Value,
			Type:     fieldType,
			Editable: true,
		})
	}
	return fields
}

func (cm *ConfigManager) UpdateConfigFromForm(config *Config, formValues map[string][]string) (*Config, error) {
	var newFields []ConfigField

	for _, field := range config.Fields {
		values, exists := formValues[field.Name]
		formValue := ""
		if exists && len(values) > 0 {
			formValue = values[0]
		}

		var newValue interface{}
		if formValue != "" {
			convertedValue, err := convertValue(formValue, field.Value)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования значения для поля %s: %w", field.Name, err)
			}
			newValue = convertedValue
		} else {
			if reflect.TypeOf(field.Value).Kind() == reflect.Bool {
				newValue = false
			} else {
				newValue = field.Value
			}
		}

		// Special handling for boolean fields (checkboxes)
		if reflect.TypeOf(field.Value).Kind() == reflect.Bool {
			if exists && len(values) > 0 && values[0] == "on" {
				newValue = true
			} else {
				newValue = false
			}
		}

		newFields = append(newFields, ConfigField{
			Name:  field.Name,
			Value: newValue,
		})
	}

	return &Config{Fields: newFields}, nil
}

func (cm *ConfigManager) readConfigFromFile() (*Config, error) {
	file, err := os.Open(cm.configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Try to parse as new format first
	var newConfig Config
	decoder := json.NewDecoder(file)
	decoder.UseNumber()

	if err := decoder.Decode(&newConfig); err != nil {
		// If new format fails, try old format for backward compatibility
		file.Seek(0, 0) // Reset file pointer
		var oldConfig map[string]interface{}
		if err := decoder.Decode(&oldConfig); err != nil {
			return nil, err
		}

		// Convert old format to new format
		return cm.convertOldFormatToNew(oldConfig), nil
	}

	// Process json.Number in new format
	for i, field := range newConfig.Fields {
		if num, ok := field.Value.(json.Number); ok {
			if iVal, err := num.Int64(); err == nil {
				newConfig.Fields[i].Value = iVal
			} else if fVal, err := num.Float64(); err == nil {
				newConfig.Fields[i].Value = fVal
			}
		}
	}

	return &newConfig, nil
}

func (cm *ConfigManager) convertOldFormatToNew(oldConfig map[string]interface{}) *Config {
	// Define field order for old format migration
	fieldOrder := []string{"fieldText", "intData", "boolValue"}

	var fields []ConfigField
	for _, fieldName := range fieldOrder {
		if value, exists := oldConfig[fieldName]; exists {
			// Process json.Number
			if num, ok := value.(json.Number); ok {
				if iVal, err := num.Int64(); err == nil {
					value = iVal
				} else if fVal, err := num.Float64(); err == nil {
					value = fVal
				}
			}
			fields = append(fields, ConfigField{
				Name:  fieldName,
				Value: value,
			})
		}
	}

	// Add any additional fields not in the predefined order
	for name, value := range oldConfig {
		found := false
		for _, fieldName := range fieldOrder {
			if name == fieldName {
				found = true
				break
			}
		}
		if !found {
			// Process json.Number
			if num, ok := value.(json.Number); ok {
				if iVal, err := num.Int64(); err == nil {
					value = iVal
				} else if fVal, err := num.Float64(); err == nil {
					value = fVal
				}
			}
			fields = append(fields, ConfigField{
				Name:  name,
				Value: value,
			})
		}
	}

	return &Config{Fields: fields}
}

func (cm *ConfigManager) saveConfigToFile(config *Config) error {
	file, err := os.Create(cm.configFile)
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
