package fetcher

import (
	"io"
	"net/http"
	"net/url"
)

//Preset fetch preset.
type Preset []Command

//Clone clone preset.
func (p *Preset) Clone() *Preset {
	cmds := make([]Command, len(*p))
	copy(cmds, *p)
	ep := BuildPreset(cmds...)
	return ep
}

//With clone preset with commands.
func (p *Preset) With(cmds ...Command) *Preset {
	preset := BuildPreset(append(*p, cmds...)...)
	return preset
}

//Append clone and append preset with given presets in order.
func (p *Preset) Append(presets ...*Preset) *Preset {
	var cmds = p.Commands()
	for k := range presets {
		cmds = append(cmds, presets[k].Commands()...)
	}
	return BuildPreset(cmds...)
}

//Commands return preset commands.
func (p *Preset) Commands() []Command {
	return []Command(*p)
}

//EndPoint create new preset with given pathprefix and method.
func (p *Preset) EndPoint(pathprefix string, method string) *Preset {
	return p.With(PathPrefix(pathprefix), Method(method))
}

//Fetch fetch request.
//Preset and commands will exec on new fetcher by which fetching response.
//Return http response and any error if raised.
func (p *Preset) Fetch(cmds ...Command) (*Response, error) {
	return Fetch(p.With(cmds...).Commands()...)
}

//FetchWithBody fetch request with given body.
//Return http response and any error if raised.
func (p *Preset) FetchWithBody(body io.Reader) (*Response, error) {
	return p.Fetch(Body(body))
}

//FetchAndParse fetch request and prase response with given parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func (p *Preset) FetchAndParse(preset Parser) (*Response, error) {
	return FetchAndParse(p, preset)
}

//FetchWithBodyAndParse fetch request and prase response with given preset ,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func (p *Preset) FetchWithBodyAndParse(body io.Reader, preset Parser) (*Response, error) {
	return FetchWithBodyAndParse(p, body, preset)
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

//ServerInfo server info struct
type ServerInfo struct {
	//URL server host url
	URL string
	//Header http header
	Header http.Header
	//Method http method
	Method string
}

//CreatePreset create new preset.
//Return preset created and any error raised.
func (s *ServerInfo) CreatePreset() (*Preset, error) {
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}
	p := BuildPreset(URL(u), Method(s.Method), Header(s.Header))
	return p, nil
}

//Server http server config struct
type Server struct {
	ServerInfo
	Client Client
}

//CreatePreset create new preset.
//Return preset created and any error raised.
func (s *Server) CreatePreset() (*Preset, error) {
	var err error

	doer, err := s.Client.CreateDoer()
	if err != nil {
		return nil, err
	}
	p, err := s.ServerInfo.CreatePreset()
	if err != nil {
		return nil, err
	}
	return p.With(SetDoer(doer)), nil
}

//PresetFactory preset factory.
type PresetFactory interface {
	//CreatePreset create new preset.
	//Return preset created and any error raised.
	CreatePreset() (*Preset, error)
}

func MustPreset(f PresetFactory) *Preset {
	p, err := f.CreatePreset()
	if err != nil {
		panic(err)
	}
	return p
}
