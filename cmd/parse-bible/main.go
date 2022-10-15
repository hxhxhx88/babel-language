package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"babel/openapi/gen/babelapi"
)

func main() {
	app := &cli.App{
		Name:  "parse-bible",
		Usage: "parse bible downloaded from https://www.ph4.org/b4_mobi.php",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "path to the downloaded xml file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "lang",
				Aliases:  []string{"l"},
				Usage:    "language ISO639-3 code",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "corpus",
				Aliases:  []string{"c"},
				Usage:    "id of the curpus to add the translation",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "server",
				Aliases:  []string{"s"},
				Usage:    "address to the babel server",
				Required: true,
			},
		},
		Action: run,
	}

	must(app.Run(os.Args))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Document struct {
	XMLName xml.Name `xml:"XMLBIBLE"`
	Books   []struct {
		Number   int    `xml:"bnumber,attr"`
		Name     string `xml:"bname,attr"`
		Chapters []struct {
			Number int `xml:"cnumber,attr"`
			Verses []struct {
				Number  int    `xml:"vnumber,attr"`
				Content string `xml:",chardata"`
			} `xml:"VERS"`
		} `xml:"CHAPTER"`
	} `xml:"BIBLEBOOK"`
}

func run(ctx *cli.Context) error {
	filePath := ctx.String("path")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return errors.WithStack(err)
	}

	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return errors.WithStack(err)
	}

	blocks, err := convert(&doc)
	if err != nil {
		return errors.WithStack(err)
	}

	tranlation := &babelapi.TranslationDraft{
		LanguageIso6393: ctx.String("lang"),
		Blocks:          &blocks,
	}

	id, err := create(ctx, tranlation)
	if err != nil {
		return err
	}
	log.Printf("created translation with id %s", id)

	return nil
}

func convert(doc *Document) ([]babelapi.BlockDraft, error) {
	uuids := make(map[string]int)

	var blocks []babelapi.BlockDraft
	for _, book := range doc.Books {
		buuid := fmt.Sprintf("book.%d", book.Number)
		if uuids[buuid] > 0 {
			return nil, errors.Errorf("duplicate uuid: %s", buuid)
		}
		uuids[buuid] = len(blocks)

		blocks = append(blocks, babelapi.BlockDraft{
			Content: strings.TrimSpace(book.Name),
			Rank:    1,
			Uuid:    buuid,
		})
		for _, chapter := range book.Chapters {
			cuuid := fmt.Sprintf("%s/chapter.%d", buuid, chapter.Number)
			if uuids[cuuid] > 0 {
				return nil, errors.Errorf("duplicate uuid: %s", cuuid)
			}
			uuids[cuuid] = len(blocks)

			blocks = append(blocks, babelapi.BlockDraft{
				Content: fmt.Sprintf("Chapter %d", chapter.Number),
				Rank:    2,
				Uuid:    cuuid,
			})
			for _, verse := range chapter.Verses {
				vuuid := fmt.Sprintf("%s/verse.%d", cuuid, verse.Number)
				if idx := uuids[vuuid]; idx > 0 {
					blocks[idx].Content += " " + strings.TrimSpace(verse.Content)
					log.Printf("combined duplicate verses [%s] to [%s]", vuuid, blocks[idx].Content)
				} else {
					uuids[vuuid] = len(blocks)

					blocks = append(blocks, babelapi.BlockDraft{
						Content: strings.TrimSpace(verse.Content),
						Rank:    3,
						Uuid:    vuuid,
					})
				}
			}
		}
	}

	return blocks, nil
}

func create(ctx *cli.Context, translation *babelapi.TranslationDraft) (string, error) {
	client, err := babelapi.NewClientWithResponses(ctx.String("server"))
	if err != nil {
		return "", errors.WithStack(err)
	}

	corpusId := ctx.String("corpus")
	req := babelapi.CreateCorpusTranslationJSONRequestBody{
		Translation: *translation,
	}
	resp, err := client.CreateCorpusTranslationWithResponse(ctx.Context, corpusId, req)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.JSON200 == nil {
		return "", errors.Errorf("http error: %s", resp.Status())
	}

	return resp.JSON200.Id, nil
}
