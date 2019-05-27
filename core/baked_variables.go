package core

//AppVersion is set at compile-time to the git tag
var AppVersion = "development"

//AppReleaseExecutableName is set to the executable name which should be downloaded by the updater next time the app is updated
var AppReleaseExecutableName = "biedaprint-linux-amd64"

//AppReleasesEndpoint the endpoint from where to download new versions of the app
var AppReleasesEndpoint = "https://api.github.com/repos/alufers/biedaprint/releases"

var AppSingleReleaseEndpoint = "https://api.github.com/repos/alufers/biedaprint/releases/tags/%s"
