#ifndef CURA_INTERFACE
#define CURA_INTERFACE
#include <string>
#include "../CuraEngine/src/settings/Settings.h"
#include "../CuraEngine/src/Scene.h"
#include "../CuraEngine/src/Slice.h"

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

    void addSetting(const char *key, const char *value)
    {
        this->wrapped_slice->scene.settings.add(key, value);
    };

    void performSlice()
    {
       this->wrapped_slice->compute();
    };

private:
    cura::Slice *wrapped_slice;

};

#endif