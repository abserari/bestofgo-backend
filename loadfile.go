package main

import "io/ioutil"

// LoadFile -
func LoadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//  UnmarshalListProjects -
func UnmarshalListProjects(data []byte) ([]OriginProject, error) {
	var OriginProjects []OriginProject

	err := json.Unmarshal(data, &OriginProjects)
	if err != nil {
		return nil, err
	}

	return OriginProjects, nil
}

//  UnmarshalListHeroes -
func UnmarshalListHeroes(data []byte) ([]OriginHero, error) {
	var OriginHeroes []OriginHero

	err := json.Unmarshal(data, &OriginHeroes)
	if err != nil {
		return nil, err
	}

	return OriginHeroes, nil
}
