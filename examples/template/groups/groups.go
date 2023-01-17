package groups

type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ParentIDs   []string `json:"parent_ids,omitempty"`
	ParentNames []string `json:"parent_names,omitempty"`
}
