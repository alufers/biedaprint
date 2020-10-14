#ifndef CURA_INTERFACE
#define CURA_INTERFACE
#include <string>
#include <iostream>
#include "../CuraEngine/src/settings/Settings.h"
#include "../CuraEngine/src/Scene.h"
#include "../CuraEngine/src/Slice.h"
#include "../CuraEngine/src/utils/floatpoint.h"

class CuraInterface
{
public:
    CuraInterface()
    {
        this->wrapped_slice = new cura::Slice(1);
    };
    const char *getAllSettingsString()
    {
        return this->wrapped_slice->scene.settings.getAllSettingsString().c_str();
    };

    void addGlobalSetting(const char *key, const char *value)
    {
        this->wrapped_slice->scene.settings.add(key, value);
    };

    void addExtruder()
    {
        this->wrapped_slice->scene.extruders.emplace_back(this->wrapped_slice->scene.extruders.size(), &this->wrapped_slice->scene.settings);
    };

    int loadModelIntoMeshGroup(const char *path)
    {

        const cura::FMatrix3x3 transformation = this->wrapped_slice->scene.settings.get<cura::FMatrix3x3>("mesh_rotation_matrix"); //The transformation applied to the model when loaded.
        std::cout << "NO SEGFAULT YEAAAAH";
        if (!loadMeshIntoMeshGroup(&this->wrapped_slice->scene.mesh_groups[0], path, transformation, this->wrapped_slice->scene.extruders[0].settings))
        {
            return errno;
        }
        return 0;
    };
    void performSlice()
    {
        this->wrapped_slice->compute();
    };

private:
    cura::Slice *wrapped_slice;
};

#endif