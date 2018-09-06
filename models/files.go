package models

// --------------------------------------
//
//  Bitbucket API files types
//
// --------------------------------------

// Files structure of FileList
type Files []string

type FilesID struct {
	Files      Files
	ProjectKey string
	RepoSlug   string
}
