package envloader

import (
	"log"
	"os"
	"strings"
)

func Load() {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		log.Fatalf("error parsing file: %v", err.Error())
	}

	lines := strings.Split(string(bytes), "\n")

	for i, line := range lines {
		if len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "#") {
			arr := strings.Split(line, "=")
			if len(arr) != 2 {
				log.Fatalf("invalid format at line %v\n%v", i+1, line)
			}

			key, value := arr[0], arr[1]
			err := os.Setenv(key, value)
			if err != nil {
				return
			}
		}
	}
}
