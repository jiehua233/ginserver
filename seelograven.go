package main

import (
	"errors"
	log "github.com/cihub/seelog"
	"github.com/getsentry/raven-go"
	"strings"
)

type RavenReciver struct {
	Client *raven.Client
}

func (r *RavenReciver) ReceiveMessage(message string, level log.LogLevel, ctx log.LogContextInterface) error {
	if r.Client == nil {
		return nil
	}
	// Gin Server的log独立
	if strings.HasPrefix(message, "Gin") {
		return nil
	}
	// 只收集Error以上级别的错误
	if level < log.ErrorLvl {
		return nil
	}
	// StackTrace 构造
	trace := raven.NewStacktrace(0, 2, nil)
	fpack, fname := FunctionNameByFunc(ctx.Func())
	frame := &raven.StacktraceFrame{
		Filename:     ctx.FullPath(),
		Function:     fname,
		Module:       fpack,
		AbsolutePath: ctx.FullPath(),
		Lineno:       ctx.Line(),
		InApp:        true,
	}
	// TODO Is Here a Bug ?
	// trace.Frames = append(trace.Frames[:0], append([]*raven.StacktraceFrame{frame}, trace.Frames[0:]...)...)
	trace.Frames = append([]*raven.StacktraceFrame{frame}, trace.Frames...)
	packet := raven.NewPacket(message, raven.NewException(errors.New(message), trace))
	r.Client.Capture(packet, nil)

	return nil
}

func (r *RavenReciver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}

func (r *RavenReciver) Flush() {

}

func (r *RavenReciver) Close() error {
	return nil
}

//

func FunctionNameByFunc(f string) (pack, name string) {
	// We get this:
	//runtime/debug.*T·ptrmethod
	// and want this:
	//  pack = runtime/debug
	//name = *T.ptrmethod
	name = f
	if idx := strings.LastIndex(name, "."); idx != -1 {
		pack = name[:idx]
		name = name[idx+1:]
	}
	name = strings.Replace(name, "·", ".", -1)
	return
}
