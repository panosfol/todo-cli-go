package util

import (
	"log"
	"bufio"
	"os"
)

func Scanner (s *string) *string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	*s = scanner.Text()
	return s
}
