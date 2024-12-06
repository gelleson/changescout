package validators

import "testing"

func TestIsValidURL(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid url with http",
			args: args{
				u: "http://www.example.com",
			},
			want: true,
		},
		{
			name: "valid url with https",
			args: args{
				u: "https://www.google.com",
			},
			want: true,
		},
		{
			name: "invalid url with missing scheme",
			args: args{
				u: "www.missing-scheme.com",
			},
			want: false,
		},
		{
			name: "invalid url with missing host",
			args: args{
				u: "https://",
			},
			want: false,
		},
		{
			name: "empty url",
			args: args{
				u: "",
			},
			want: false,
		},
		{
			name: "url with unexpected characters",
			args: args{
				u: "://invalid#url",
			},
			want: false,
		},
		{
			name: "url with only scheme",
			args: args{
				u: "http://",
			},
			want: false,
		},
		{
			name: "file url scheme",
			args: args{
				u: "file:///path/to/file",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.args.u); got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
