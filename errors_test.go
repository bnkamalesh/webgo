package webgo

import (
	"testing"
)

func Test_loggerWithCfg(t *testing.T) {
	t.Parallel()
	cfgs := []logCfg{
		LogCfgDisableDebug,
		LogCfgDisableInfo,
		LogCfgDisableWarn,
		LogCfgDisableError,
		LogCfgDisableFatal,
	}
	l := loggerWithCfg(nil, nil, cfgs...)
	if l.debug != nil {
		t.Errorf("expected debug to be nil, got %v", l.debug)
	}
	if l.err != nil {
		t.Errorf("expected err to be nil, got %v", l.err)
	}
	if l.fatal != nil {
		t.Errorf("expected fatal to be nil, got %v", l.fatal)
	}
	if l.info != nil {
		t.Errorf("expected info to be nil, got %v", l.info)
	}
	if l.warn != nil {
		t.Errorf("expected warn to be nil, got %v", l.warn)
	}
}
