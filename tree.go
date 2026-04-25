package sml

import "fmt"

const (
	PROCPARVALUETAGVALUE       = 0x01
	PROCPARVALUETAGPERIODENTRY = 0x02
	PROCPARVALUETAGTUPELENTRY  = 0x03
	PROCPARVALUETAGTIME        = 0x04
)

type TreePath []OctetString

type Tree struct {
	ParameterName  OctetString
	ParameterValue *ProcParValue
	ChildList      []*Tree
}

type ProcParValue struct {
	Tag         uint8
	Value       *Value
	PeriodEntry *PeriodEntry
	TupelEntry  *TupelEntry
	Time        *Time
}

type PeriodEntry struct {
	ObjName        OctetString
	Unit           uint8
	Scaler         int8
	Value          Value
	ValueSignature OctetString
}

type TupelEntry struct {
	ServerID        OctetString
	SecIndex        Time
	Status          uint64
	UnitPA          uint8
	ScalerPA        int8
	ValuePA         int64
	UnitR1          uint8
	ScalerR1        int8
	ValueR1         int64
	UnitR4          uint8
	ScalerR4        int8
	ValueR4         int64
	SignaturePAR1R4 OctetString
	UnitMA          uint8
	ScalerMA        int8
	ValueMA         int64
	UnitR2          uint8
	ScalerR2        int8
	ValueR2         int64
	UnitR3          uint8
	ScalerR3        int8
	ValueR3         int64
	SignatureMAR2R3 OctetString
}

func TreePathParse(buf *Buffer) (TreePath, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	if err := ExpectType(buf, TYPELIST); err != nil {
		return nil, err
	}

	elems := BufGetNextLength(buf)
	path := make(TreePath, 0, elems)

	for i := 0; i < elems; i++ {
		entry, err := OctetStringParse(buf)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			path = append(path, entry)
		}
	}

	return path, nil
}

func TreeParse(buf *Buffer) (*Tree, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	tree := &Tree{}
	var err error

	if err := Expect(buf, TYPELIST, 3); err != nil {
		return nil, err
	}

	if tree.ParameterName, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	if tree.ParameterValue, err = ProcParValueParse(buf); err != nil {
		return nil, err
	}

	if !BufOptionalIsSkipped(buf) {
		if err := ExpectType(buf, TYPELIST); err != nil {
			return nil, err
		}

		elems := BufGetNextLength(buf)
		tree.ChildList = make([]*Tree, 0, elems)

		for i := 0; i < elems; i++ {
			child, err := TreeParse(buf)
			if err != nil {
				return nil, err
			}
			if child != nil {
				tree.ChildList = append(tree.ChildList, child)
			}
		}
	}

	return tree, nil
}

func ProcParValueParse(buf *Buffer) (*ProcParValue, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	ppv := &ProcParValue{}

	if err := Expect(buf, TYPELIST, 2); err != nil {
		return nil, err
	}

	tag, err := U8Parse(buf)
	if err != nil {
		return nil, err
	}
	ppv.Tag = tag

	switch ppv.Tag {
	case PROCPARVALUETAGVALUE:
		v, err := ValueParse(buf)
		if err != nil {
			return nil, err
		}
		ppv.Value = &v
	case PROCPARVALUETAGPERIODENTRY:
		period, err := PeriodEntryParse(buf)
		if err != nil {
			return nil, err
		}
		ppv.PeriodEntry = period
	case PROCPARVALUETAGTUPELENTRY:
		tupel, err := TupelEntryParse(buf)
		if err != nil {
			return nil, err
		}
		ppv.TupelEntry = tupel
	case PROCPARVALUETAGTIME:
		time, err := TimeParse(buf)
		if err != nil {
			return nil, err
		}
		ppv.Time = &time
	default:
		return nil, fmt.Errorf("Invalid proc parameter value tag %02x", ppv.Tag)
	}

	return ppv, nil
}

func TupelEntryParse(buf *Buffer) (*TupelEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	tupel := &TupelEntry{}
	var err error

	if err := Expect(buf, TYPELIST, 23); err != nil {
		return nil, err
	}

	if tupel.ServerID, err = OctetStringParse(buf); err != nil {
		return nil, err
	}
	if tupel.SecIndex, err = TimeParse(buf); err != nil {
		return nil, err
	}
	if tupel.Status, err = U64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitPA, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerPA, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValuePA, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitR1, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerR1, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValueR1, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitR4, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerR4, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValueR4, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.SignaturePAR1R4, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitMA, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerMA, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValueMA, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitR2, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerR2, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValueR2, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.UnitR3, err = U8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ScalerR3, err = I8Parse(buf); err != nil {
		return nil, err
	}
	if tupel.ValueR3, err = I64Parse(buf); err != nil {
		return nil, err
	}

	if tupel.SignatureMAR2R3, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	return tupel, nil
}

func PeriodEntryParse(buf *Buffer) (*PeriodEntry, error) {
	if BufOptionalIsSkipped(buf) {
		return nil, nil
	}

	period := &PeriodEntry{}
	var err error

	if err := Expect(buf, TYPELIST, 5); err != nil {
		return nil, err
	}

	if period.ObjName, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	if period.Unit, err = U8Parse(buf); err != nil {
		return nil, err
	}

	if period.Scaler, err = I8Parse(buf); err != nil {
		return nil, err
	}

	if period.Value, err = ValueParse(buf); err != nil {
		return nil, err
	}

	if period.ValueSignature, err = OctetStringParse(buf); err != nil {
		return nil, err
	}

	return period, nil
}
