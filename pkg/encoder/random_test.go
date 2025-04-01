package encoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRandomSeed = 42
const testAlphabet = "ABCDEFGHabcdefgh12340_"

func TestRandom_Encode(t *testing.T) {
	type args struct {
		text   string
		length int
	}

	testCases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				text:   "https://youtube.com/watch?v=",
				length: 10,
			},
			want: "dbe1df_g1f",
		},
		{
			name: "zero length",
			args: args{
				text:   "https://youtube.com/watch?v=",
				length: 0,
			},
			want: "",
		},
		{
			name: "negative length",
			args: args{
				text:   "https://youtube.com/watch?v=",
				length: -1,
			},
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRandom(testAlphabet, testRandomSeed)

			value := r.Encode(tc.args.text, tc.args.length)
			assert.Equal(t, tc.want, value)
		})
	}
}
