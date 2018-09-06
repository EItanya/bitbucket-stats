package models

// --------------------------------------
//
//  Bitbucket API project types
//
// --------------------------------------

// Project JSON form of single project
type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
}

// FilterProjects method to filter saved projects
func FilterProjects(data *[]Project, projects []string) []Project {
	if len(projects) == 0 {
		return *data
	}
	filteredProjects := make([]Project, 0)
	ch := make(chan []Project)
	for _, val := range projects {
		go filterProjects(val, data, ch)
	}
	for range projects {
		filteredProjects = append(filteredProjects, <-ch...)
	}
	return filteredProjects
}

func filterProjects(val string, p *[]Project, ch chan []Project) {
	pm := make([]Project, 0)
	for _, v := range *p {
		if v.Key == val || string(v.ID) == val {
			pm = append(pm, v)
		}
	}
	ch <- pm
}