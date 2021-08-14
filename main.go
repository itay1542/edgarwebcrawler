package edgarwebcrawler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// downloader := edgarwebcrawler.NewIdxDownloader()

	// err := downloader.Download(".\\storage\\indices", "https://www.sec.gov/Archives/edgar/full-index", true, 2006)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	uris := ".\\storage\\submission_uris.txt"
	uriPrefix := "https://www.sec.gov/Archives/"
	err := filepath.Walk(".\\storage\\indices\\full-index", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			parser, err := edgarwebcrawler.NewIdxReader(file)
			fmt.Println(path)
			if err != nil {
				return err
			}
			saveFiles(parser, uris, uriPrefix)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}

func saveFiles(parser edgarwebcrawler.IdxReader, destFile, addPrefix string) error {
	var row *edgarwebcrawler.IdxRow
	//skip the annoying ------ row
	_, err := parser.ReadRow()
	hasSeenFormFour := false
	f, err := os.OpenFile(destFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for {
		row, err = parser.ReadRow()
		if err != nil {
			log.Fatal(err)
		}
		if hasSeenFormFour && row.FormType != "4" {
			break
		}
		if row.FormType == "4" {
			hasSeenFormFour = true
			if _, err := f.WriteString(fmt.Sprintf("%s\n", addPrefix+row.FileName)); err == nil {
				log.Printf("Successfully appended file name")
			}

		}
		fmt.Printf("%+v \n", row)
	}
	return nil
}
