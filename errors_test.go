package webgo

import (
	"bytes"
	"strings"
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

func TestLogging(t *testing.T) {
	logmsg := "hello world"
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	logh := loggerWithCfg(stdout, stderr)

	logh.Debug(logmsg)
	stdStr := stdout.String()
	outstr := strings.TrimSpace(stdStr)
	outstr = outstr[(len(outstr))-len(logmsg):]
	if outstr != logmsg {
		t.Fatalf(
			"expected output '%s', got '%s'",
			logmsg,
			outstr,
		)
	}
	if stdStr[0:5] != "Debug" {
		t.Fatalf(
			"expected '%s' in the beginning of log, got '%s'",
			"Debug",
			stdStr[0:6],
		)
	}

	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	logh = loggerWithCfg(stdout, stderr)
	logh.Info(logmsg)
	stdStr = stdout.String()
	outstr = strings.TrimSpace(stdStr)
	outstr = outstr[(len(outstr))-len(logmsg):]
	if outstr != logmsg {
		t.Fatalf(
			"expected output '%s', got '%s'",
			logmsg,
			outstr,
		)
	}
	if stdStr[0:4] != "Info" {
		t.Fatalf(
			"expected '%s' in the beginning of log, got '%s'",
			"Info",
			stdStr[0:4],
		)
	}

	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	logh = loggerWithCfg(stdout, stderr)
	logh.Warn(logmsg)
	stdStr = stderr.String()
	outstr = strings.TrimSpace(stdStr)
	outstr = outstr[(len(outstr))-len(logmsg):]
	if outstr != logmsg {
		t.Fatalf(
			"expected output '%s', got '%s'",
			logmsg,
			outstr,
		)
	}
	if stdStr[0:4] != "Warn" {
		t.Fatalf(
			"expected '%s' in the beginning of log, got '%s'",
			"Warn",
			stdStr[0:4],
		)
	}

	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	logh = loggerWithCfg(stdout, stderr)
	logh.Error(logmsg)
	stdStr = stderr.String()
	outstr = strings.TrimSpace(stdStr)
	outstr = outstr[(len(outstr))-len(logmsg):]
	if outstr != logmsg {
		t.Fatalf(
			"expected output '%s', got '%s'",
			logmsg,
			outstr,
		)
	}
	if stdStr[0:5] != "Error" {
		t.Fatalf(
			"expected '%s' in the beginning of log, got '%s'",
			"Error",
			stdStr[0:4],
		)
	}

}
