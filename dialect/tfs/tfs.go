package tfs

import (
	"io/ioutil"
	"log"

	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/castro/dialect"
)

// TFS struct that defines a dialect for
// the forgotten server
type TFS struct {
	stages []dialect.Stage
}

// Name shows the dialect name
func (t TFS) Name() string {
	return "the forgotten server"
}

// Version shows the dialect version
func (t TFS) Version() string {
	return "0.1 alpha-preview"
}

// LoadStages loads server xml stages
func (t *TFS) LoadStages() error {
	file, err := ioutil.ReadFile(util.Config.Datapack + "/data/xml/stages.xml")
	if err != nil {
		return err
	}
	log.Println(file)
	return nil
}

func (t TFS) GetStages() []dialect.Stage {
	return t.stages
}
