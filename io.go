package edgarwebcrawler

import "bufio"

func ReadLine(reader *bufio.Reader) (string, error) {
	line := ""
	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	for b != byte('\n') {
		line += string(b)
		b, err = reader.ReadByte()
		if err != nil {
			return "", err
		}
	}
	return line, nil
}