package backend

import (
	"errors"

	"babel/app/storage"
)

type Options struct {
	storageCorpus storage.Corpus
}

func (o *Options) Validate() error {
	if o.storageCorpus == nil {
		return errors.New("missing corpus storage")
	}
	return nil
}

type Option func(*Options)

func WithCorpusStorage(s storage.Corpus) Option {
	return func(o *Options) {
		o.storageCorpus = s
	}
}
