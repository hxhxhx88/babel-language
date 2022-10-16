package backend

import (
	"errors"

	"babel/app/storage"
)

type Options struct {
	storageCorpus      storage.Corpus
	storageTranslation storage.Translation
	storageBlock       storage.Block
}

func (o *Options) Validate() error {
	if o.storageCorpus == nil {
		return errors.New("missing corpus storage")
	}
	if o.storageTranslation == nil {
		return errors.New("missing translation storage")
	}
	if o.storageBlock == nil {
		return errors.New("missing block storage")
	}
	return nil
}

type Option func(*Options)

func WithCorpusStorage(s storage.Corpus) Option {
	return func(o *Options) {
		o.storageCorpus = s
	}
}

func WithTranslationStorage(s storage.Translation) Option {
	return func(o *Options) {
		o.storageTranslation = s
	}
}

func WithBlockStorage(s storage.Block) Option {
	return func(o *Options) {
		o.storageBlock = s
	}
}
