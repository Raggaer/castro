package dialect

import (
	"io/ioutil"

	"github.com/raggaer/castro/app/util"
)

// TFS struct that defines a dialect for
// the forgotten server
type TFS struct {
	stages []Stage
}

// Registers the TFS dialect
func init() {
	List.Register("tfs", &TFS{})
}

// Name shows the dialect name
func (t TFS) Name() string {
	return "the forgotten server"
}

// Version shows the dialect version
func (t TFS) Version() string {
	return "1.2"
}

// Start executes all the needed stuff for the dialect
// to work correctly
func (t TFS) Start() error {
	return nil
}

// LoadStages loads server xml stages
func (t *TFS) LoadStages() error {
	_, err := ioutil.ReadFile(util.Config.Datapack + "/data/xml/stages.xml")
	if err != nil {
		return err
	}
	//log.Println(file)
	return nil
}

func (t TFS) GetStages() []Stage {
	return t.stages
}
