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
			name: "valid url",
			args: args{
				u: "https://www.google.com",
			},
			want: true,
		},
		{
			name: "invalid url",
			args: args{
				u: "///invalid",
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
