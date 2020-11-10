package wxpay

import (
	"encoding/xml"
	"io"
)

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

//XMLMap ...
type XMLMap map[string]string

//UnmarshalXML ...
func (m *XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = XMLMap{}
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
func (m XMLMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "xml",
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

// CData ...
type CData string

// UnmarshalXML ...
func (c *CData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	cdata := struct{
		Data string `xml:",cdata"`
	}{}

	if err := d.DecodeElement(&cdata, &start); err != nil {
		return err
	}
	*c = CData(cdata.Data)
	return nil
}

// MarshalXML ...
func (c *CData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	cdata := struct{
		Data string `xml:",cdata"`
	}{Data: string(*c)}
	if err := e.EncodeElement(cdata, start); err != nil {
		return err
	}

	return nil
}
