package domain

import (
	"path"
	"strings"
)

type ArchiveName struct {
	value string
}

func (a *ArchiveName) String() string {
	return a.value
}

func NewArchiveName(src SourcePath, paths []EntryPath) (ArchiveName, error) {
	if len(paths) < 1 {
		return ArchiveName{}, ErrEmptySelection
	}

	if len(paths) == 1 {
		return ArchiveName{baseOr(src, paths[0].FS()) + ".zip"}, nil
	}
	return ArchiveName{baseOr(src, commonParent(paths)) + ".zip"}, nil
}

func commonParent(paths []EntryPath) string {
	common := strings.Split(path.Dir(paths[0].FS()), "/")
	for _, p := range paths[1:] {
		parts := strings.Split(path.Dir(p.FS()), "/")
		n := 0
		for n < len(common) && n < len(parts) && common[n] == parts[n] {
			n++
		}
		common = common[:n]
	}
	return strings.Join(common, "/")
}

// baseOr returns the last segment of p, or the source name for the root.
func baseOr(src SourcePath, p string) string {
	if b := path.Base(p); b != "." && b != "/" {
		return b
	}
	return src.String()
}
