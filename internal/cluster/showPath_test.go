package cluster

import "testing"

func Test_GetFilePath(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				fileName: "../../config/test.tpl",
			},
			want: "/Users/onko/go/src/github.com/kubernauts/tk8/config/test.tpl",
		}, {
			name: "test2",
			args: args{
				fileName: "../../configs/templates/test.tpl",
			},
			want: "/Users/onko/go/src/github.com/kubernauts/tk8/configs/templates/test.tpl",
		}, {
			name: "test3",
			args: args{
				fileName: "../../test.tpl",
			},
			want: "/Users/onko/go/src/github.com/kubernauts/tk8/test.tpl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilePath(tt.args.fileName); got != tt.want {
				t.Errorf("GetFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
