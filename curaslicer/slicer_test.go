package curaslicer

import (
	"log"
	"testing"
)

func TestSlicerRuns(t *testing.T) {
	iface := NewCuraInterface()
	defer DeleteCuraInterface(iface)
	iface.AddSetting("name", "val")
	iface.PerformSlice()
	log.Print(iface.GetAllSettingsString())
}
