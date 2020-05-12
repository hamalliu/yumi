package wxpay

import (
	"encoding/xml"
	"io"
)

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type XmlMap map[string]string

func (m *XmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = XmlMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

// MarshalXML allows type H to be used with xml.Marshal.
func (m XmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range m {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
