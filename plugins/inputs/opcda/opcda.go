package opcda

// opcda.go

import (
	"fmt"
	"log"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"

	"github.com/konimarti/opc"
)

type TagSettings struct {
	ItemID    string     `toml:"item"`
	FieldName string     `toml:"name"`
	TagsSlice [][]string `toml:"tags"`
}

type OpcTag struct {
	fieldName string
	tags      map[string]string
}

type Opcda struct {
	// Configuration
	MeasurementName string        `toml:"name"`
	Server          string        `toml:"server"`
	Nodes           []string      `toml:"nodes"`
	OpcTagConf      []TagSettings `toml:"items"`

	Log telegraf.Logger `toml:"-"`

	// internal values
	client     opc.Connection
	opcItemIds []string
	items      map[string]OpcTag
}

func (s *Opcda) Description() string {
	return "Retrieve data from OPC DA servers"
}

func (s *Opcda) SampleConfig() string {
	return `
## Measurement name
# name = "sim"

## OPC DA server
# server = "OI.SIM.1"

## First successful node in the list is used
# nodes = ["localhost"]

## Node ID configuration
## item        - OPC DA Item Name or ItemID
## name        - field name to use in the output (optional)
## tags        - extra tags to be added to the output metric (optional)
## Example:
# items = [
#   { item = "PORT.PLC.int" },
#   { item = "PORT.PLC.float", name = "randomFloat", tags = [
#     [
#       "equipment",
#       "maker1",
#     ],
#     [
#       "input",
#       "temperature",
#     ],
#   ] },
# ]
`
}

func tagsSliceToMap(tags [][]string) (map[string]string, error) {
	m := make(map[string]string)
	for i, tag := range tags {
		if len(tag) != 2 {
			return nil, fmt.Errorf("tag %d needs 2 values, has %d: %v", i+1, len(tag), tag)
		}
		if tag[0] == "" {
			return nil, fmt.Errorf("tag %d has empty name", i+1)
		}
		if tag[1] == "" {
			return nil, fmt.Errorf("tag %d has empty value", i+1)
		}
		if _, ok := m[tag[0]]; ok {
			return nil, fmt.Errorf("tag %d has duplicate key: %v", i+1, tag[0])
		}
		m[tag[0]] = tag[1]
	}
	return m, nil
}

// InitNodes is read config and init OPC tags
func (s *Opcda) InitTags() error {
	s.items = make(map[string]OpcTag)

	for _, t := range s.OpcTagConf {
		tags, err := tagsSliceToMap(t.TagsSlice)
		if err != nil {
			log.Fatal(err)
		}

		s.items[t.ItemID] = OpcTag{
			fieldName: t.FieldName,
			tags:      tags,
		}

		s.opcItemIds = append(s.opcItemIds, t.ItemID)
	}

	return nil
}

// Init is for setup, and validating config.
func (s *Opcda) Init() error {
	return s.InitTags()
}

func (s *Opcda) Start(acc telegraf.Accumulator) error {
	s.Log.Info("Connecting to ", s.Server, s.Nodes)

	var err error
	s.client, err = opc.NewConnection(
		s.Server,
		s.Nodes,
		s.opcItemIds,
	)

	return err
}

func (s *Opcda) Stop() {
	s.client.Close()
}

func (s *Opcda) Gather(acc telegraf.Accumulator) error {
	items := s.client.Read()

	for tagName, opcItem := range items {
		if opcItem.Quality == opc.OPCQualityGood {
			var fieldname string
			opcTag, ok := s.items[tagName]

			if ok && len(opcTag.fieldName) > 0 {
				fieldname = opcTag.fieldName
			} else {
				fieldname = tagName
			}

			fields := make(map[string]interface{})
			fields[fieldname] = opcItem.Value
			acc.AddFields(s.MeasurementName, fields, opcTag.tags, opcItem.Timestamp)
		}
	}

	return nil
}

func init() {
	inputs.Add("opcda", func() telegraf.Input { return &Opcda{} })
}
