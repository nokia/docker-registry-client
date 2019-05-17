package registry

import (
	"fmt"
	"net/http"
)

type tagsResponse struct {
	Tags []string `json:"tags"`
}

func (registry *Registry) Tags(repository string) (tags []string, err error) {
	url := registry.url("/v2/%s/tags/list", repository)

	var response tagsResponse
	for {
		registry.Logf("registry.tags url=%s repository=%s", url, repository)
		url, err = registry.getPaginatedJson(url, &response)
		switch err {
		case ErrNoMorePages:
			tags = append(tags, response.Tags...)
			return tags, nil
		case nil:
			tags = append(tags, response.Tags...)
			continue
		default:
			return nil, err
		}
	}
}

// DeleteTag is Hub specific to delete a tag from the docker registry.
func (registry *Registry) DeleteTag(repository, tag string) (err error) {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s/", repository, tag)

	registry.Logf("registry.tag.delete url=%s repository=%s tag=%s", url, repository, tag)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := registry.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	return nil
}
