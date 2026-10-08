package api

// ProjectConfig is the part of the Config the project declares itself: api.Config
// embeds it, so every field typed here is read as sandbox.Config.<Field>,
// beside the ProjectName and Version the build renders from project.yaml:
//
//	type ProjectConfig struct {
//		Port int
//	}
//
// Fill it in sandbox/constructors/config/constructor.go, after NewConfig — that
// file is written once and is yours too. A field is written in the builtin
// types or in a type of this package, like every other of sandbox/api/.
//
// Written once by `agnos start` and then yours: no build rewrites
// this file.
type ProjectConfig struct {
}
