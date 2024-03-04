package session

import (
	"encoding/hex"

	"google.golang.org/protobuf/proto"

	"github.com/andatoshiki/toshiki-rssbot/internal/log"
)

// Marshal encodes to a string
func Marshal(a *Attachment) string {
	bytes, err := proto.Marshal(a)
	if err != nil {
		log.Errorf("marshal attachment failed, %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

// Unmarshal Attachment parses the transmitted information from a string
func UnmarshalAttachment(data string) (*Attachment, error) {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	a := &Attachment{}
	if err := proto.Unmarshal(bytes, a); err != nil {
		return nil, err
	}
	return a, nil
}
