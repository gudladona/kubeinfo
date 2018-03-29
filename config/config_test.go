package config

import "testing"

func TestReadConfigFromFile(t *testing.T) {
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestConfigReadSuccess",
			args: args{
				cfgFile: "./testdata/test-config.cfg",
			},
			wantErr: false,
		},
		{
			name: "TestConfigReadFail",
			args: args{
				cfgFile: "./somecfg.cfg",
			},
			wantErr: true,
		},
		{
			name: "TestBadConfig",
			args: args{
				cfgFile: "./testdata/bad-config.cfg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadConfigFromFile(tt.args.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListenPort(t *testing.T) {
	// Init
	ReadConfigFromFile("./testdata/test-config.cfg")

	tests := []struct {
		name string
		want int
	}{
		{
			name: "TestConfigReadSuccess",
			want: 8080,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListenPort(); got != tt.want {
				t.Errorf("ListenPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
