package snippet

import (
	"encoding/json"
	"fmt"
	stacktrace "github.com/palantir/stacktrace"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	util "github.com/da-moon/cli-snippets/internal/util"
	color "github.com/fatih/color"
)

// SnippetsMeta ...
type SnippetsMeta struct {
	Snippets    []*jsonSnippet `json:"snippets"`
	IsMetaDirty bool           `json:"is_meta_dirty"`
	fileLoc     string
	snippetsDir string
}

type jsonSnippet struct {
	FileLoc string `json:"file_loc"`
	Title   string `json:"title"`
}

// SetFileLoc ...
func (sm *SnippetsMeta) SetFileLoc(fileLoc string) {
	sm.fileLoc = fileLoc
}

// SetSnippetsDir ...
func (sm *SnippetsMeta) SetSnippetsDir(path string) {
	sm.snippetsDir = path
}

// SyncWithSnippets ...
func (sm *SnippetsMeta) SyncWithSnippets() error {
	for _, s := range sm.Snippets {
		// update file location to always use snippetsDir
		s.FileLoc = path.Join(sm.snippetsDir, getSnippetFileName(s.Title))
		snippet, err := LoadSnippet(s.FileLoc)
		if err != nil {
			return err
		}
		// if title changed in snippet file, then update both file name and title in meta
		if s.Title != snippet.Title {
			newFileName := getSnippetFileName(snippet.Title)
			newFilePath := path.Join(sm.snippetsDir, newFileName)
			s.Title = snippet.Title
			if err = os.Rename(s.FileLoc, newFilePath); err != nil {
				return err
			}
			s.FileLoc = newFilePath
		}
	}
	sm.IsMetaDirty = false
	if err := sm.Save(); err != nil {
		return err
	}
	return nil
}

// Save ...
func (sm *SnippetsMeta) Save() error {
	if _, err := os.Stat(sm.fileLoc); os.IsNotExist(err) {
		return err
	}
	data, err := json.MarshalIndent(sm, util.JSONMarshalPrefix, util.JSONMarshalIndent)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(sm.fileLoc, data, 0644); err != nil {
		return err
	}
	return nil
}

// Save new snippet into snippetsDir and update snippets meta file
func (sm *SnippetsMeta) SaveNewSnippet(snippet *Snippet) error {
	// check for duplicate
	if sm.isDuplicate(snippet.Title) {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		newTitle := fmt.Sprintf("%s-%s", snippet.Title, t)
		color.Red("Snippet with title \"%s\" already existed - saving as \"%s\"", snippet.Title, newTitle)
		snippet.Title = newTitle
	}
	// save snippet file
	if err := snippet.Save(sm.snippetsDir); err != nil {
		return err
	}
	// save to snippets meta file
	jsonSnippet := &jsonSnippet{
		Title:   snippet.Title,
		FileLoc: snippet.fileLoc,
	}
	sm.Snippets = append(sm.Snippets, jsonSnippet)
	if err := sm.Save(); err != nil {
		return err
	}
	return nil
}

func (sm *SnippetsMeta) isDuplicate(title string) bool {
	for _, s := range sm.Snippets {
		if s.Title == title {
			return true
		}
	}
	return false
}

// DeleteSnippet ...
func (sm *SnippetsMeta) DeleteSnippet(title string) error {
	idx, err := sm.findJsonSnippetIndex(title)
	if err != nil {
		return err
	}
	s := sm.Snippets[idx]
	fmt.Printf("Deleting snippet %s... ", s.Title)
	// delete snippet file
	fileLoc := s.FileLoc
	if err = os.Remove(fileLoc); err != nil {
		color.Red("Failure")
		return err
	}
	// delete from snippets meta
	sm.Snippets = append(sm.Snippets[:idx], sm.Snippets[idx+1:]...)
	if err = sm.Save(); err != nil {
		color.Red("Failure")
		return err
	}
	color.Green("Success")
	return nil
}

// FindSnippet ...
func (sm *SnippetsMeta) FindSnippet(title string) (*Snippet, error) {
	idx, err := sm.findJsonSnippetIndex(title)
	if err != nil {
		return nil, err
	}
	s, err := LoadSnippet(sm.Snippets[idx].FileLoc)
	if err != nil {
		if os.IsNotExist(err) {
			err = stacktrace.Propagate(err, "Snippets directory path: %s, check if snippet file exists inside.", sm.snippetsDir)
			return nil, err
		}
		return nil, err
	}
	return s, nil
}

func (sm *SnippetsMeta) findJsonSnippetIndex(title string) (int, error) {
	idx := -1
	for i, snp := range sm.Snippets {
		if snp.Title == title {
			idx = i
			break
		}
	}
	if idx == -1 {
		err := stacktrace.NewError("Could not find snippet with name: %s", title)
		return idx, err
	}
	return idx, nil
}

// GetSnippetTitles ...
func (sm *SnippetsMeta) GetSnippetTitles() []string {
	titles := make([]string, len(sm.Snippets))
	for idx, s := range sm.Snippets {
		titles[idx] = s.Title
	}
	return titles
}
