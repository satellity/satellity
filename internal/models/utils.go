package models

import (
	"crypto/md5"
	"io"

	"github.com/gofrs/uuid"
)

func generateUniqueID(args ...string) string {
	h := md5.New()
	for _, arg := range args {
		io.WriteString(h, arg)
	}
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	id, err := uuid.FromBytes(sum)
	if err != nil {
		panic(err)
	}
	return id.String()
}
