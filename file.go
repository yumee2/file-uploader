package main

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}
