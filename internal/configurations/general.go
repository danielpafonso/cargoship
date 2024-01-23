package configurations

import (
	"fmt"
	"strings"
	"time"
)

// replaceDatePlaceholder substitutes datetime format string from configurations with current date time
func replaceDatePlaceholder(filename string) string {
	start := strings.Index(filename, "{")
	end := strings.Index(filename, "}")
	if start != -1 && end != -1 {
		dateFormat := filename[start+1 : end]
		return fmt.Sprintf(
			"%s%s%s",
			filename[:start],
			time.Now().UTC().Format(dateFormat),
			filename[end+1:],
		)
	}
	return filename
}
