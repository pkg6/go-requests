package requests

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestBase64Decode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Decode(tt.args.str)
			if err != nil != tt.wantErr {
				t.Errorf("Base64Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Encode(tt.args.str); got != tt.want {
				t.Errorf("Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64File(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64File(tt.args.path)
			if err != nil != tt.wantErr {
				t.Errorf("Base64File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64File() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Reader(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Reader(tt.args.reader)
			if err != nil != tt.wantErr {
				t.Errorf("Base64Reader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64Reader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.str); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5File(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "test1",
			path: "./.testdata/github.png",
			want: "404a5db8eec868e5f29b9d20b0395094",
		},
		{
			name: "test2",
			path: "./.testdata/github.png",
			want: "404a5db8eec868e5f29b9d20b0395094",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Md5File(tt.path)
			if got != tt.want {
				t.Errorf("Md5File() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5Reader(t *testing.T) {
	f, _ := os.Open("./.testdata/github.png")
	tests := []struct {
		name string
		file *os.File
		want string
	}{
		{
			name: "test1",
			file: f,
			want: "404a5db8eec868e5f29b9d20b0395094",
		},
		{
			name: "test2",
			file: f,
			want: "404a5db8eec868e5f29b9d20b0395094",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = tt.file.Seek(0, 0)
			got, _ := Md5Reader(tt.file)
			if got != tt.want {
				t.Errorf("Md5Reader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha1(tt.args.str); got != tt.want {
				t.Errorf("Sha1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha256(tt.args.str); got != tt.want {
				t.Errorf("Sha256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLDecode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLDecode(tt.args.str)
			if err != nil != tt.wantErr {
				t.Errorf("URLDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLEncode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLEncode(tt.args.str); got != tt.want {
				t.Errorf("URLEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64StdEncoding(t *testing.T) {
	type args struct {
		base64Str string
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64StdEncoding(tt.args.base64Str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Base64StdEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}
