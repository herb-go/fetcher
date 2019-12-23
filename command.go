package fetcher

//Command fetch command interface which used to modify fetch
type Command interface {
	//Exec exec command to modify fetcher.
	//Return any error if raised.
	Exec(*Fetcher) error
}

//CommandFunc command func
type CommandFunc func(*Fetcher) error

//Exec exec command to modify fetcher.
//Return any error if raised.
func (f CommandFunc) Exec(e *Fetcher) error {
	return f(e)
}

//Exec exec given commands to fetcher by order.
//Return any error if raised
func Exec(f *Fetcher, b ...Command) error {
	var err error
	for k := range b {
		err = b[k].Exec(f)
		if err != nil {
			return err
		}
	}
	return nil
}
