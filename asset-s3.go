package kasset

type S3Asset struct {
	key, secret, bucket string
}

func NewS3(key, secret, bucket string) (*S3Asset, error) {
	s3 := new(S3Asset)
	s3.key = key
	s3.secret = secret
	s3.bucket = bucket
	return s3, nil
}

func (sfs *S3Asset) Save(name string, bs []byte) error {
	return nil
}

func (sfs *S3Asset) Read(name string) ([]byte, error) {
	return []byte{}, e
}

func (sfs *S3Asset) Delete(name string) error {
	return nil
}
