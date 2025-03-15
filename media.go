package genanki

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)


type MediaFile struct {
	Filename string
	Data     []byte
	Hash     string
}


func NewMediaFile(filename string, data []byte) *MediaFile {
	return &MediaFile{
		Filename: SanitizeFilename(filename),
		Data:     data,
		Hash:     GenerateMediaHash(data),
	}
}


func NewMediaFileFromPath(path string) (*MediaFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read media file: %v", err)
	}

	filename := filepath.Base(path)
	return NewMediaFile(filename, data), nil
}


func NewMediaFileFromReader(filename string, reader io.Reader) (*MediaFile, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read media data: %v", err)
	}

	return NewMediaFile(filename, data), nil
}


func (m *MediaFile) GetMimeType() string {
	ext := strings.ToLower(filepath.Ext(m.Filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		
		return "application/octet-stream"
	}
	return mimeType
}


func (m *MediaFile) IsImage() bool {
	mimeType := m.GetMimeType()
	return strings.HasPrefix(mimeType, "image/")
}


func (m *MediaFile) IsAudio() bool {
	mimeType := m.GetMimeType()
	return strings.HasPrefix(mimeType, "audio/")
}


func (m *MediaFile) IsVideo() bool {
	mimeType := m.GetMimeType()
	return strings.HasPrefix(mimeType, "video/")
}


func (p *Package) AddMediaFromPath(path string) error {
	mediaFile, err := NewMediaFileFromPath(path)
	if err != nil {
		return err
	}
	p.media[mediaFile.Filename] = mediaFile.Data
	return nil
}


func (p *Package) AddMediaFromReader(filename string, reader io.Reader) error {
	mediaFile, err := NewMediaFileFromReader(filename, reader)
	if err != nil {
		return err
	}
	p.media[mediaFile.Filename] = mediaFile.Data
	return nil
}


func (p *Package) GetMediaFiles() []*MediaFile {
	files := make([]*MediaFile, 0, len(p.media))
	for filename, data := range p.media {
		files = append(files, NewMediaFile(filename, data))
	}
	return files
}


func (p *Package) GetMediaFile(filename string) *MediaFile {
	if data, ok := p.media[filename]; ok {
		return NewMediaFile(filename, data)
	}
	return nil
}


func (p *Package) RemoveMedia(filename string) {
	delete(p.media, filename)
}


func (p *Package) ClearMedia() {
	p.media = make(map[string][]byte)
}


func (p *Package) GetMediaCount() int {
	return len(p.media)
}


func (p *Package) GetMediaSize() int64 {
	var total int64
	for _, data := range p.media {
		total += int64(len(data))
	}
	return total
}
