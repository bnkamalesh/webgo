package webgo

import (
	"bytes"
	"testing"
)

func TestGlobalLoggerConfig(t *testing.T) {
	const logMsg = "hello world"
	type args struct {
		cfgs []logCfg
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
	}{
		{
			name: "disable all",
			args: args{
				cfgs: []logCfg{
					LogCfgDisableDebug,
					LogCfgDisableInfo,
					LogCfgDisableWarn,
					LogCfgDisableError,
					LogCfgDisableFatal,
				},
			},
			wantStdout: "",
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		GlobalLoggerConfig(stdout, stderr, tt.args.cfgs...)

		LOGHANDLER.Debug(logMsg)
		LOGHANDLER.Info(logMsg)
		LOGHANDLER.Warn(logMsg)
		LOGHANDLER.Error(logMsg)
		LOGHANDLER.Fatal(logMsg)

		if gotStdout := stdout.String(); gotStdout != tt.wantStdout {
			t.Errorf("GlobalLoggerConfig() = %v, want %v", gotStdout, tt.wantStdout)
		}
		if gotStderr := stderr.String(); gotStderr != tt.wantStderr {
			t.Errorf("GlobalLoggerConfig() = %v, want %v", gotStderr, tt.wantStderr)
		}
	}
}
