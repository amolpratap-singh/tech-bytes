package maps

import "errors"

type Dictionary map[string]string

var ErrNotFound = errors.New("could not find the word you were looking for")


func Search(dict map[string]string, word string) string {
	return dict[word]
}

func (d Dictionary) CustomSearch(word string) (string, error) {
	defination, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	return defination, nil
}

func (d Dictionary) Add(word, defination string) {
	d[word] = defination
}