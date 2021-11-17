package phc

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

const (
	fieldsDelimiter = "$"
	versionParam    = "v"
	kvDelimiter     = "="
	paramsDelimiter = ","
)

var (
	ErrInvalidEncoding = errors.New("invalid encoding")
	B64                = base64.RawStdEncoding
)

// https://github.com/P-H-C/phc-string-format
type PHC struct {
	AlgID   string
	Version *uint32
	Params  map[string]string
	Salt    []byte
	Hash    []byte

	UseUpperVersion bool
}

func (p *PHC) String() string {
	if p.AlgID == "" {
		panic("PHC: AlgID is empty")
	}
	var b strings.Builder
	b.WriteString(fieldsDelimiter + p.AlgID)
	if p.Version != nil {
		var s string
		if p.UseUpperVersion {
			s = strings.ToUpper(versionParam)
		} else {
			s = versionParam
		}
		b.WriteString(
			fieldsDelimiter + s + kvDelimiter + strconv.FormatUint(uint64(*p.Version), 10))
	}
	if p.Params != nil && len(p.Params) != 0 {
		b.WriteString(fieldsDelimiter)
		i := 0
		for k, v := range p.Params {
			b.WriteString(k + kvDelimiter + v)
			if i != len(p.Params)-1 {
				b.WriteString(paramsDelimiter)
			}
			i++
		}
	}
	if p.Salt != nil && len(p.Salt) != 0 {
		b.WriteString(fieldsDelimiter + B64.EncodeToString(p.Salt))
		if p.Hash != nil && len(p.Hash) != 0 {
			b.WriteString(fieldsDelimiter + B64.EncodeToString(p.Hash))
		}
	}
	return b.String()
}

func FromString(str string) (*PHC, error) {
	l := strings.Split(str, fieldsDelimiter)
	pos := 0
	if len(l) < 2 || len(l) > 6 || l[pos] != "" {
		return nil, ErrInvalidEncoding
	}
	if pos++; l[pos] == "" {
		return nil, ErrInvalidEncoding
	}
	res := &PHC{AlgID: l[pos]}
	if pos++; len(l) > pos &&
		(strings.HasPrefix(l[pos], strings.ToUpper(versionParam)+kvDelimiter) ||
			(strings.HasPrefix(l[pos], versionParam+kvDelimiter) &&
				!strings.Contains(l[pos], paramsDelimiter))) {

		i, err := strconv.ParseUint(l[pos][len(versionParam)+len(kvDelimiter):], 10, 32)
		if err != nil {
			return nil, ErrInvalidEncoding
		}
		j := uint32(i)
		res.Version = &j
		pos++
	}
	if len(l) > pos && strings.Contains(l[pos], kvDelimiter) {
		kvs := strings.Split(l[pos], paramsDelimiter)
		res.Params = make(map[string]string, len(kvs))
		for i := range kvs {
			kv := strings.Split(kvs[i], kvDelimiter)
			if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
				return nil, ErrInvalidEncoding
			}
			res.Params[kv[0]] = kv[1]
		}
		pos++
	}
	if len(l) > pos && l[pos] != "" {
		var err error
		res.Salt, err = B64.DecodeString(l[pos])
		if err != nil {
			return nil, ErrInvalidEncoding
		}
		if pos++; len(l) > pos && l[pos] != "" {
			res.Hash, err = B64.DecodeString(l[pos])
			if err != nil {
				return nil, ErrInvalidEncoding
			}
			pos++
		}
	}
	if len(l) > pos {
		return nil, ErrInvalidEncoding
	}
	return res, nil
}
