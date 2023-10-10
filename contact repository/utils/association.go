package utils

import "strings"


func GetAssociations(includes []string,associations []string)[]string{
	matchedAssociations := []string{}
	if(len(includes)==0 || len(associations)==0){
		return matchedAssociations
	}
	for _,i := range includes {
		for _,a := range associations {
			if strings.Contains(a,i){
				matchedAssociations = append(matchedAssociations, a)
			}
		}
	}
	return matchedAssociations
}