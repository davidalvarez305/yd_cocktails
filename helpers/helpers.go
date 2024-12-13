package helpers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/davidalvarez305/yd_cocktails/types"
)

func ServeDynamicPartialTemplate(w http.ResponseWriter, ctx types.DynamicPartialTemplate) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	template, err := BuildStringFromTemplate(ctx.TemplatePath, ctx.TemplateName, ctx.Data)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error parsing template.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(template))
}

func InsertHTMLIntoEmailTemplate(templatePath, templateName, emailBody string, data any) (string, error) {
	// Read the wrapper template file.
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}

	// Parse the wrapper template.
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}

	// Define the inline template with the dynamic email body content.
	tmpl, err = tmpl.New("content.html").Parse(emailBody)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}

	// Render the template to a string using a strings.Builder.
	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}

	return output.String(), nil
}

func BuildStringFromTemplate(templatePath, templateName string, data any) (string, error) {
	var output string
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	var body strings.Builder
	err = tmpl.Execute(&body, data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	return body.String(), err
}

func ServeContent(w http.ResponseWriter, templateFilePaths []string, data any) error {
	tmpl, err := template.ParseFiles(templateFilePaths...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building template.", http.StatusInternalServerError)
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error executing template.", http.StatusInternalServerError)
		return err
	}

	return nil
}

func BuildFile(outputPath string, templateFilePaths []string, data any) error {
	tmpl, err := template.ParseFiles(templateFilePaths...)

	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(outputPath)

	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)

	if err != nil {
		return err
	}

	err = f.Close()

	if err != nil {
		return err
	}

	return err
}

func GetUserIPFromRequest(r *http.Request) string {
	// Try to get the IP from the X-Forwarded-For header (used by proxies and load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Try to get the IP from the X-Real-IP header (used by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to using the RemoteAddr from the request
	ip := r.RemoteAddr
	return ip
}

func RemoveCountryCode(phoneNumber string) string {
	if strings.HasPrefix(phoneNumber, "+1") {
		return phoneNumber[2:]
	}
	return phoneNumber
}

func GenerateSequence(start, end int) []int {
	var sequence []int
	for i := start; i <= end; i++ {
		sequence = append(sequence, i)
	}
	return sequence
}

func CalculateMaxPages(totalRows, itemsPerPage int) int {
	if totalRows <= 0 {
		return 1
	}
	return (totalRows + itemsPerPage - 1) / itemsPerPage
}

func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func SafeStringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetMarketingCookiesFromRequestOrForm(r *http.Request, cookieName, formName string) *string {
	cookie, err := r.Cookie(cookieName)
	if err == nil && cookie.Value != "" {
		return &cookie.Value
	}

	return GetStringPointerFromForm(r, formName)
}

func GetStringPointerFromForm(r *http.Request, key string) *string {
	if values, ok := r.Form[key]; ok && len(values) > 0 {
		return &values[0]
	}
	return nil
}

func SafeStringToIntPointer(value string) *int {
	if val, err := strconv.Atoi(value); err == nil {
		return &val
	}
	return nil
}

func SafeStringToInt64Pointer(value string) *int64 {
	if val, err := strconv.ParseInt(value, 10, 64); err == nil {
		return &val
	}
	return nil
}

func GetIntPointerFromForm(r *http.Request, key string) *int {
	if values, ok := r.Form[key]; ok && len(values) > 0 {
		if i, err := strconv.Atoi(values[0]); err == nil {
			return &i
		}
	}
	return nil
}

func GetInt64PointerFromForm(r *http.Request, key string) *int64 {
	if values, ok := r.Form[key]; ok && len(values) > 0 {
		if i, err := strconv.ParseInt(values[0], 10, 64); err == nil {
			return &i
		}
	}
	return nil
}

func GetFloat64PointerFromForm(r *http.Request, key string) *float64 {
	if values, ok := r.Form[key]; ok && len(values) > 0 {
		if i, err := strconv.ParseFloat(values[0], 64); err == nil {
			return &i
		}
	}
	return nil
}

func GetBoolPointerFromForm(r *http.Request, key string) *bool {
	if values, ok := r.Form[key]; ok && len(values) > 0 {
		parsedValue, err := strconv.ParseBool(values[0])
		if err == nil {
			return &parsedValue
		}
	}
	return nil
}

func ParseInt64(value string) int64 {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsedValue
}

func GetFirstIDAfterPrefix(r *http.Request, prefix string) (int, error) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, prefix)
	trimmedPath = strings.Trim(trimmedPath, "/")

	parts := strings.Split(trimmedPath, "/")

	for _, part := range parts {
		id, err := strconv.Atoi(part)
		if err == nil {
			return id, nil
		}
	}

	return 0, fmt.Errorf("no valid ID found in path")
}

func GetSecondIDFromPath(r *http.Request, prefix string) (int, error) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, prefix)
	parts := strings.Split(trimmedPath, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path format")
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid  id: %w", err)
	}

	return id, nil
}

func GetThirdIDFromPath(r *http.Request, prefix string) (int, error) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, prefix)
	parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path format")
	}

	fmt.Println(parts)

	idStr := parts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}

func IsNumeric(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func IsMobileRequest(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	mobileKeywords := []string{"Mobile", "Android", "iPhone", "iPad", "iPod"}

	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}
	return false
}

func GetStringValueFromField(field reflect.Value) string {
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		return field.Elem().String()
	} else if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}
