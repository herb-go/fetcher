package fetcher

import (
	"io"
	"net/http"
	"path"
)

//Preset fetch preset.
type Preset struct {
	prev    *Preset
	command Command
}

//Exec exec command to modify fetcher.
//Return any error if raised.
func (p *Preset) Exec(f *Fetcher) error {
	if p == nil {
		return nil
	}
	err := p.prev.Exec(f)
	if err != nil {
		return err
	}
	return p.command.Exec(f)
}

// CloneWith clone preset with commands.
// Alias for Concat
func (p *Preset) CloneWith(cmds ...Command) *Preset {
	return p.Concat(cmds...)
}

// With clone preset with commands.
// Alias for Concat
func (p *Preset) With(cmds ...Command) *Preset {
	return p.Concat(cmds...)
}
func (p *Preset) concatCommand(cmd Command) *Preset {
	return &Preset{
		prev:    p,
		command: cmd,
	}
}

//Concat concat preset with given commands
func (p *Preset) Concat(cmds ...Command) *Preset {
	preset := p
	for _, v := range cmds {
		preset = preset.concatCommand(v)
	}
	return preset
}
func (p *Preset) appendPreset(preset *Preset) *Preset {

	if preset == nil {
		return p
	}
	newpreset := p.appendPreset(preset.prev)
	return &Preset{
		prev:    newpreset,
		command: preset.command,
	}
}

//Append clone and append preset with given presets in order.
func (p *Preset) Append(presets ...*Preset) *Preset {

	preset := p
	for _, v := range presets {
		preset = preset.appendPreset(v)
	}
	return preset
}

//Commands return preset commands.
func (p *Preset) Commands() []Command {
	if p == nil {
		return []Command{}
	}
	cmds := p.prev.Commands()
	cmds = append(cmds, p.command)
	return cmds
}

//EndPoint create new preset with given suffix and method.
func (p *Preset) EndPoint(method string, suffix string) *Preset {
	return p.Concat(PathSuffix(suffix), Method(method))
}

//Fetch fetch request.
//Preset and commands will exec on new fetcher by which fetching response.
//Return http response and any error if raised.
func (p *Preset) Fetch(cmds ...Command) (*Response, error) {
	return Fetch(p.Concat(cmds...))
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

//FetchWithJSONBodyAndParse fetch request and prase response with given preset ,body as json and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func (p *Preset) FetchWithJSONBodyAndParse(body interface{}, preset Parser) (*Response, error) {
	return FetchAndParse(p.Concat(JSONBody(body)), preset)
}

//NewPreset create new preset
func NewPreset() *Preset {
	return nil
}

//Concat create preset with given commands.
func Concat(cmds ...Command) *Preset {
	return NewPreset().Concat(cmds...)
}

//BuildPreset build new preset with given commands
func BuildPreset(cmds ...Command) *Preset {
	return NewPreset().Concat(cmds...)
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

//MergeURL clone and merge with given url
func (s *ServerInfo) MergeURL(url string) *ServerInfo {
	si := s.Clone()
	si.URL = url
	return si
}

//Join clone and join with given urlpath
func (s *ServerInfo) Join(urlpath string) *ServerInfo {
	si := s.Clone()
	si.URL = path.Join(si.URL, urlpath)
	return si
}

//CreatePreset create new preset.
//Return preset created and any error raised.
func (s *ServerInfo) CreatePreset() (*Preset, error) {
	p := BuildPreset(URL(s.URL), Method(s.Method), Header(s.Header))
	return p, nil
}

//IsEmpty check if server info is empty
//IsEmpty will return true if equal nil or s.URL is empty string
func (s *ServerInfo) IsEmpty() bool {
	return s == nil || s.URL == ""
}

//Clone clone a new serverinfo
func (s *ServerInfo) Clone() *ServerInfo {
	si := &ServerInfo{
		URL:    s.URL,
		Header: s.Header.Clone(),
		Method: s.Method,
	}
	return si
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
	return p.Concat(SetDoer(doer)), nil
}

//Clone clone a new server config
func (s *Server) Clone() *Server {
	return &Server{
		ServerInfo: *s.ServerInfo.Clone(),
		Client:     *s.Client.Clone(),
	}
}

//MergeURL clone and merge with given url
func (s *Server) MergeURL(url string) *Server {
	return &Server{
		ServerInfo: *s.ServerInfo.MergeURL(url),
		Client:     *s.Client.Clone(),
	}
}

//Join clone and join with given urlpath
func (s *Server) Join(urlpath string) *Server {
	return &Server{
		ServerInfo: *s.ServerInfo.Join(urlpath),
		Client:     *s.Client.Clone(),
	}
}

//PresetFactory preset factory.
type PresetFactory interface {
	//CreatePreset create new preset.
	//Return preset created and any error raised.
	CreatePreset() (*Preset, error)
}

//MustPreset create preset by given preset factory.
//Return preset created.
//Panic if any error raised.
func MustPreset(f PresetFactory) *Preset {
	p, err := f.CreatePreset()
	if err != nil {
		panic(err)
	}
	return p
}
