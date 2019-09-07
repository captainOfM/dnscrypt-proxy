package main

import (
	"strings"

	"github.com/jedisct1/dlog"
	"github.com/miekg/dns"
)

type PluginFirefox struct {
}

func (plugin *PluginFirefox) Name() string {
	return "firefox"
}

func (plugin *PluginFirefox) Description() string {
	return "Work around Firefox taking over DNS"
}

func (plugin *PluginFirefox) Init(proxy *Proxy) error {
	dlog.Noticef("Firefox workaround initialized")
	return nil
}

func (plugin *PluginFirefox) Drop() error {
	return nil
}

func (plugin *PluginFirefox) Reload() error {
	return nil
}

func (plugin *PluginFirefox) Eval(pluginsState *PluginsState, msg *dns.Msg) error {
	questions := msg.Question
	if len(questions) != 1 {
		return nil
	}
	question := questions[0]
	if question.Qclass != dns.ClassINET || (question.Qtype != dns.TypeA && question.Qtype != dns.TypeAAAA) {
		return nil
	}
	qName := strings.ToLower(question.Name)
	if qName != "use-application-dns.net." && !strings.HasSuffix(qName, ".use-application-dns.net.") {
		return nil
	}
	synth, err := EmptyResponseFromMessage(msg)
	if err != nil {
		return err
	}
	synth.Rcode = dns.RcodeNameError
	pluginsState.synthResponse = synth
	pluginsState.action = PluginsActionSynth
	pluginsState.returnCode = PluginsReturnCodeSynth
	return nil
}
