package postgres

import (
	"babel/openapi/gen/babelapi"
	"strconv"
)

type mCorpusRecord struct {
	Id                      int    `db:"id"`
	Title                   string `db:"title"`
	OriginalLanguageIso6393 string `db:"original_language_iso_639_3"`
}

func (rec *mCorpusRecord) ToOpenApi() *babelapi.Corpus {
	return &babelapi.Corpus{
		Id:                      strconv.Itoa(rec.Id),
		Title:                   rec.Title,
		OriginalLanguageIso6393: rec.OriginalLanguageIso6393,
	}
}

type mTranslationRecord struct {
	Id              int    `db:"id"`
	CorpusId        int    `db:"corpus_id"`
	Title           string `db:"title"`
	LanguageIso6393 string `db:"language_iso_639_3"`
}

func (rec *mTranslationRecord) ToOpenApi() *babelapi.Translation {
	return &babelapi.Translation{
		Id:              strconv.Itoa(rec.Id),
		CorpusId:        strconv.Itoa(rec.CorpusId),
		Title:           rec.Title,
		LanguageIso6393: rec.LanguageIso6393,
	}
}

type mBlockRecord struct {
	Id            int    `db:"id"`
	TranslationId int    `db:"corpus_id"`
	Content       string `db:"content"`
	Rank          int    `db:"rank"`
	Uuid          string `db:"uuid"`
}

func (rec *mBlockRecord) ToOpenApi() *babelapi.Block {
	return &babelapi.Block{
		Id:            strconv.Itoa(rec.Id),
		TranslationId: strconv.Itoa(rec.TranslationId),
		Content:       rec.Content,
		Rank:          rec.Rank,
		Uuid:          rec.Uuid,
	}
}
