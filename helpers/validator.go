package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// TranslateErrorMessage menangani error validasi dan database menjadi map yang lebih ramah
func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Handle validasi dari validator.v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = "Invalid email format"
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field)
			default:
				errorsMap[field] = "Invalid value"
			}
		}
	}

	// Handle GORM error: Duplicate entry
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			field := extractDuplicateField(err.Error())
			if field != "" {
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			} else {
				errorsMap["Error"] = "Duplicate entry"
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["Error"] = "Record not found"
		}
	}

	return errorsMap
}

// IsDuplicateEntryError mengecek apakah error adalah duplicate entry
func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}

// extractDuplicateField mencoba mengekstrak nama kolom dari error "Duplicate entry"
func extractDuplicateField(errMsg string) string {
	// Contoh error MySQL: Error 1062: Duplicate entry 'test@example.com' for key 'users.email'
	// Kita ambil bagian setelah 'for key' lalu extract nama field
	re := regexp.MustCompile(`for key '(\w+\.)?(\w+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) == 3 {
		// Hasilkan kapitalisasi nama field # strings.Title deprecated
		// return strings.Title(matches[2])

		// Title-case  untuk Unicode
		c := cases.Title(language.English)
		return c.String(matches[2])
	}
	return ""
}
