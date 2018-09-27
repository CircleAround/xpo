package stdkit

import (
	"bufio"
	"io"
)

func FindLine(reader *bufio.Reader, callback func(line string) bool) (bool, error) {
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
		// fmt.Println(string(line))
		if callback(string(line)) {
			return true, nil
		}
	}
}
