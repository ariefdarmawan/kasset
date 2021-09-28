package kasset

import "github.com/h2non/filetype"

func getFileType(buffer []byte) (string, string, error) {
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", "", err
	}
	return kind.MIME.Value, kind.Extension, nil
}

func GetFileType(buffer []byte) (string, string, error) {
	data := buffer
	if len(data) > 512 {
		data = data[:512]
	}
	return getFileType(data)
}
