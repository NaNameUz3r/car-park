package files

import (
	"car-park/telegram-bot/lib/er"
	"car-park/telegram-bot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

const (
	defaultPermissions = 0774
	saveErr            = "can't save in file"
)

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) SaveCredentials(creds *storage.Credentials) error {
	fPath := filepath.Join(s.basePath, creds.Username)
	if err := os.MkdirAll(fPath, defaultPermissions); err != nil {
		er.Wrap(saveErr, err)
	}

	fPath = filepath.Join(fPath, "credentials")

	file, err := os.Create(fPath)
	if err != nil {
		er.Wrap(saveErr, err)
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(creds); err != nil {
		er.Wrap(saveErr, err)
	}

	return nil
}

func (s Storage) RemoveCredentials(creds *storage.Credentials) error {
	fPath := filepath.Join(s.basePath, creds.Username, "credentials")

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fPath)
		return er.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExistsCredentials(userName string) (bool, error) {
	fPath := filepath.Join(s.basePath, userName, "credentials")

	switch _, err := os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("there is no credentials file %s", fPath)
		return false, er.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) GetCredentials(userName string) (*storage.Credentials, error) {
	fPath := filepath.Join(s.basePath, userName, "credentials")
	f, err := os.Open(fPath)
	if err != nil {
		return nil, er.Wrap("can't open credentials file", err)
	}

	defer func() { _ = f.Close() }()

	var creds storage.Credentials
	if err := gob.NewDecoder(f).Decode(&creds); err != nil {
		return nil, er.Wrap("cant't decode credentials file", err)
	}

	return &creds, nil
}
