package fetcher

import (
	"net/http"
	"net/url"
)

//Preset fetch preset
type Preset []Command

//Clone clone  preset
func (p *Preset) Clone() *Preset {
	cmds := make([]Command, len(*p))
	copy(cmds, *p)
	ep := BuildPreset(cmds...)
	return ep
}

//With clone preset with commands
func (p *Preset) With(cmds ...Command) *Preset {
	preset := BuildPreset(append(*p, cmds...)...)
	return preset
}

//Commands return preset commands
func (p *Preset) Commands() []Command {
	return []Command(*p)
}

//EndPoint create new preset with given pathprefix and method.
func (p *Preset) EndPoint(pathprefix string, method string) *Preset {
	return p.With(PathPrefix(pathprefix), Method(method))
}

//Fetch fetch http response.
//Preset and commands will exec on new fetcher by which fetching response.
//Return http response and any error if raised
func (p *Preset) Fetch(cmds ...Command) (*Response, error) {
	return Fetch(p.With(cmds...).Commands()...)
}

//NewPreset create new preset
func NewPreset() *Preset {
	return &Preset{}
}

//BuildPreset build new preset with given commands
func BuildPreset(cmds ...Command) *Preset {
	p := Preset(cmds)
	return &p
}

//Server http server config struct
type Server struct {
	Host   string
	Header http.Header
	Method string
	Client Client
}

//CreatePreset create new preset.
//Return preset created and any error raised.
func (s *Server) CreatePreset() (*Preset, error) {
	var err error
	u, err := url.Parse(s.Host)
	if err != nil {
		return nil, err
	}
	doer, err := s.Client.CreateDoer()
	if err != nil {
		return nil, err
	}
	p := BuildPreset(URL(u), SetDoer(doer), Method(s.Method), Header(s.Header))
	return p, nil
}

//PresetFactory preset factory.
type PresetFactory interface {
	//CreatePreset create new preset.
	//Return preset created and any error raised.
	CreateFetcher() (*Preset, error)
}
