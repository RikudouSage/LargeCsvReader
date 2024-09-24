package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CellWidget struct {
	widget.Label
	OnRightClick func(event *fyne.PointEvent)
}

func (c *CellWidget) TappedSecondary(event *fyne.PointEvent) {
	if c.OnRightClick != nil {
		c.OnRightClick(event)
	}
}

func NewCellWidget(text string, onRightClick func(*fyne.PointEvent)) *CellWidget {
	cell := &CellWidget{
		Label:        widget.Label{Text: text},
		OnRightClick: onRightClick,
	}
	cell.ExtendBaseWidget(cell)

	return cell
}
