package kanban

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"project-manager-go/data"
	"strconv"

	"github.com/google/uuid"
)

type FileManager struct {
	Path string
	HTTP string
}

func NewFileManager(path, http string) *FileManager {
	return &FileManager{
		Path: path,
		HTTP: http,
	}
}

func (fm *FileManager) GetPath() string {
	return fm.Path
}

func (fm *FileManager) GetHTTP() string {
	return fm.HTTP
}

func (fm *FileManager) GetURL(id int, name string) (string, error) {
	return url.JoinPath(fm.GetHTTP(), "uploads", strconv.Itoa(id), name)
}

func (fm *FileManager) SaveFile(file multipart.File, name string) (*data.File, error) {
	err := os.MkdirAll(fm.GetPath(), 0755)
	if err != nil {
		return nil, err
	}

	tempFile, err := os.CreateTemp(fm.GetPath(), "u*")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}

	id := int(uuid.New().ID())
	url, err := fm.GetURL(id, name)
	if err != nil {
		return nil, err
	}

	fileStore := data.File{
		ID:   id,
		Name: name,
		Path: path.Base(tempFile.Name()),
		URL:  url,
	}

	return &fileStore, nil
}

func (fm *FileManager) AddFile(dbCtx *data.DBContext, file *data.File) (err error) {
	ctx := data.NewTCtx(dbCtx)
	defer func() { err = ctx.End(err) }()

	return ctx.DB.Create(&file).Error
}

func (fm *FileManager) FindFile(w http.ResponseWriter, dbCtx *data.DBContext, id int) (err error) {
	ctx := data.NewTCtx(dbCtx)
	defer func() { err = ctx.End(err) }()

	find := data.File{}
	err = ctx.DB.Find(&find, id).Error
	if err != nil || find.ID == 0 {
		return err
	}

	file, err := os.Open(path.Join(fm.GetPath(), find.Path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}
