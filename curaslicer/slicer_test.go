package curaslicer

import (
	"log"
	"testing"
)

func TestSlicerRuns(t *testing.T) {
	iface := NewCuraInterface()
	defer DeleteCuraInterface(iface)
	iface.AddGlobalSetting("name", "val")
	iface.AddGlobalSetting("mesh_rotation_matrix", "[[1,0,0], [0,1,0], [0,0,1]]")
	iface.AddExtruder()
	log.Print(iface.GetAllSettingsString())
	errno := iface.LoadModelIntoMeshGroup("/home/alufers/Modele/db25_mounting_test.stl")
	log.Printf("errno = %v", errno)
	// iface.PerformSlice()
	log.Print(iface.GetAllSettingsString())
}
