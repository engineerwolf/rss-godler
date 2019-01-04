package main

import "testing"

func TestFileSize_UnmarshalYAML(t *testing.T) {
	type args struct {
		rawValue string
	}
	tests := []struct {
		name    string
		fs      *FileSize
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			fs:   new(FileSize),
			args: args{
				rawValue: "10Mb",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unmarshal := func(rIp interface{}) (err error) {
				testInput := tt.args.rawValue
				rIp = &testInput
				return nil
			}
			if err := tt.fs.UnmarshalYAML(unmarshal); (err != nil) != tt.wantErr {
				t.Errorf("FileSize.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
