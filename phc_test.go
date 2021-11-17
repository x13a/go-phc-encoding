package phc

import (
	"strings"
	"testing"
)

func TestPHC(t *testing.T) {
	testVectors := []string{
		"$argon2id$v=19$key=a2V5,m=4096,t=0,p=1$X1NhbHQAAAAAAAAAAAAAAA$bWh++MKN1OiFHKgIWTLvIi1iHicmHH7+Fv3K88ifFfI",
		"$scrypt$v=1$ln=15,r=8,p=1$c2FsdHNhbHQ$dGVzdHBhc3M",
		"$scrypt",
		"$scrypt$v=1",
		"$scrypt$ln=15,r=8,p=1",
		"$scrypt$c2FsdHNhbHQ",
		"$scrypt$v=1$ln=15,r=8,p=1$c2FsdHNhbHQ",
		"$scrypt$v=1$ln=15,r=8,p=1",
		"$scrypt$v=1$c2FsdHNhbHQ$dGVzdHBhc3M",
		"$scrypt$v=1$c2FsdHNhbHQ",
		"$scrypt$c2FsdHNhbHQ$dGVzdHBhc3M",
	}
	for _, s1 := range testVectors {
		v, err := FromString(s1)
		if err != nil {
			t.Errorf("%s: %s", s1, err)
			continue
		}
		s2 := v.String()
		if len(s1) != len(s2) {
			t.Errorf("%s: %s", s1, s2)
			continue
		}
		l1 := strings.Split(s1, fieldsDelimiter)
		l2 := strings.Split(s2, fieldsDelimiter)
		if len(l1) != len(l2) {
			t.Errorf("%s: %s", s1, s2)
			continue
		}
	FieldsLoop:
		for i := 0; i < len(l1); i++ {
			if l1[i] == l2[i] {
				continue
			}
			kvs1 := strings.Split(l1[i], paramsDelimiter)
			kvs2 := strings.Split(l2[i], paramsDelimiter)
			if len(kvs1) != len(kvs2) {
				t.Errorf("%s: %s", s1, s2)
				break
			}
			for _, kv1 := range kvs1 {
				found := false
				for _, kv2 := range kvs2 {
					if kv1 == kv2 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: %s", s1, s2)
					break FieldsLoop
				}
			}
		}
	}
}
