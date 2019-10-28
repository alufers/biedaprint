package core

import (
	"fmt"
)

/*
GcodeFileMetaRepositoryService handles saving the commands entered by the user in the serial console window for later use by pressing the up arrow just like in a normal shell.
*/
type GcodeFileMetaRepositoryService struct {
	app *App
}

/*
NewGcodeFileMetaRepositoryService constructs a GcodeFileMetaRepositoryService.
*/
func NewGcodeFileMetaRepositoryService(app *App) *GcodeFileMetaRepositoryService {
	return &GcodeFileMetaRepositoryService{
		app: app,
	}
}

/*
Init initializes the repository, auto-migrates the models.
*/
func (gfmrs *GcodeFileMetaRepositoryService) Init() error {
	err := gfmrs.app.DBService.DB.AutoMigrate(&GcodeFileMeta{}).Error
	if err != nil {
		return fmt.Errorf("failed to migrate GcodeFileMeta: %w", err)
	}
	return nil
}

/*
Save persists the GcodeFileMeta in the database.
*/
func (gfmrs *GcodeFileMetaRepositoryService) Save(file *GcodeFileMeta) error {
	if err := gfmrs.app.DBService.DB.Save(file).Error; err != nil {
		return fmt.Errorf("failed to save gcode file meta: %v", err)
	}
	return nil
}

/*
GetAll returns all the gcode file metas sorted by upload date in newest-to-oldest order.
*/
func (gfmrs *GcodeFileMetaRepositoryService) GetAll() (res []*GcodeFileMeta, err error) {
	if dbErr := gfmrs.app.DBService.DB.Order("upload_date DESC").Find(&res).Error; dbErr != nil {
		res = nil
		err = fmt.Errorf("failed to get all gcode file metas: %v", dbErr)
		return
	}
	return
}

/*
Delete deletes the given file from the database.
*/
func (gfmrs *GcodeFileMetaRepositoryService) Delete(file *GcodeFileMeta) error {
	if dbErr := gfmrs.app.DBService.DB.Delete(file).Error; dbErr != nil {
		return fmt.Errorf("failed to delete gcode file meta: %v", dbErr)
	}
	return nil
}

/*
GetOneByID returns a gcode file meta by id.
*/
func (gfmrs *GcodeFileMetaRepositoryService) GetOneByID(id int) (file *GcodeFileMeta, err error) {
	file = &GcodeFileMeta{}
	if dbErr := gfmrs.app.DBService.DB.Where("id = ?", id).Find(file).Error; dbErr != nil {
		file = nil
		err = fmt.Errorf("failed to find gcode file meta: %v", dbErr)
		return
	}
	return
}
