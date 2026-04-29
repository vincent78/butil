package objUtil

import "testing"

func TestSet3JsonStr(t *testing.T) {
	type args struct {
		json  string
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				json:  "",
				key:   "name1",
				value: "value1",
			},
		},
		{
			name: "test2",
			args: args{
				json:  "{\"t\":1}",
				key:   "name1",
				value: 12,
			},
		},
		{
			name: "test3",
			args: args{
				json:  "{\"t\":1}",
				key:   "t",
				value: 2,
			},
		},
		{
			name: "test4",
			args: args{
				json:  "{\"t\":1}",
				key:   "t.f.g",
				value: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Set2JsonStr(tt.args.json, tt.args.key, tt.args.value)
			if err != nil {
				t.Errorf("Set2JsonStr() error = %v", err)
			} else {
				t.Logf("got : %v", got)
			}
		})
	}
}
