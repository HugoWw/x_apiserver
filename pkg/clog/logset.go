package clog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var Logger = New("APIServer")

var GlobalClogSets = clogSets{make(map[string]*clog)}

type clogSets struct {
	logSet map[string]*clog
}

// New create a new *zap.SugaredLogger object for logging operations
// when you create a *zap.SugaredLogger object it will be registered to the global map species
func New(name string) *zap.SugaredLogger {
	atomicInfoLevel := zap.NewAtomicLevel()
	logger := initZapSugLogger(atomicInfoLevel)
	log := &clog{
		LogLevel: atomicInfoLevel,
		Logger:   logger.Named(name),
	}

	GlobalClogSets.addLog(name, log)

	return log.Logger
}

func (c *clogSets) addLog(name string, log *clog) {
	if c.logSet == nil {
		c.logSet = make(map[string]*clog)
	}
	lowerName := strings.ToLower(name)

	c.logSet[lowerName] = log
}

// GetLogAtomicLevel get the log level object for the specified module name
func (c *clogSets) GetLogAtomicLevel(name string) *zap.AtomicLevel {
	v, ok := c.logSet[name]
	if !ok {
		return nil
	}

	return &v.LogLevel
}

// SetLogAtomicLevel you can set the log level for a specific module or a list of modules
// if the module is empty, the specified log level is enabled for all log modules
func (c *clogSets) SetLogAtomicLevel(l zapcore.Level, module ...string) error {
	if len(module) > 0 {
		for _, m := range module {
			v, ok := c.logSet[m]
			if !ok {
				return fmt.Errorf("specific log module name %s not exist", m)
			}
			v.LogLevel.SetLevel(l)
		}
	} else {
		for _, v := range c.logSet {
			v.LogLevel.SetLevel(l)
		}
	}

	return nil
}
