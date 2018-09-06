package models

import "fmt"

// --------------------------------------
//
//  Bitbucket API files types
//
// --------------------------------------

// Files structure of FileList
type Files []string

func (f *Files) Unmarshal(dat interface{}) error {
	if cast, ok := dat.(*interface{}); ok {
		*cast = *f
		return nil
	}
	if cast, ok := dat.(*Files); ok {
		*cast = *f
		return nil
	}
	return fmt.Errorf("Improper type (%s) was passed into unmarshal method for Files model", dat)
}

func (f *Files) Marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case Files:
		f = &typedData
	default:
		return fmt.Errorf("Improper type (%s) was passed into marshal method for Files model", dat)
	}
	return nil
}

type FilesID struct {
	Files      Files
	ProjectKey string
	RepoSlug   string
}

func (f *FilesID) Unmarshal(dat interface{}) error {
	if cast, ok := dat.(*interface{}); ok {
		*cast = *f
		return nil
	}
	if cast, ok := dat.(*FilesID); ok {
		*cast = *f
		return nil
	}
	return fmt.Errorf("Improper type (%s) was passed into unmarshal method for FilesID model", dat)
}

func (f *FilesID) Marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case FilesID:
		f = &typedData
	default:
		return fmt.Errorf("Improper type (%s) was passed into marshal method for FilesID model", dat)
	}
	return nil
}
