package fileUtil

import "testing"

func TestPath(t *testing.T) {
	type args struct {
		rel string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				rel: "fileUtils.go",
			},
		},
		{
			name: "test2",
			args: args{
				rel: "go.mod",
			},
		},
		{
			name: "test3",
			args: args{
				rel: "conf/container/config/params.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Path(tt.args.rel)
			t.Logf("got = %v", got)
		})
	}
}
