package producthunt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty"`
	URL         string `json:"url,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

type ProductHunt struct {
	APIKey string
}

const graphqlURL = "https://api.producthunt.com/v2/api/graphql"

func fetchData(query string, apiKey string) (map[string]interface{}, error) {
	payload := map[string]string{"query": query}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (ph *ProductHunt) GetDaily() ([]Product, error) {
	query := `
    query {
      posts(order: NEWEST, first: 10) {
        edges {
          node {
            id
            name
            tagline
          }
        }
      }
    }
    `
	data, err := fetchData(query, ph.APIKey)
	if err != nil {
		return nil, err
	}

	var products []Product
	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'data'")
	}
	posts, ok := dataMap["posts"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'posts'")
	}
	edges, ok := posts["edges"].([]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'edges'")
	}

	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		product := Product{
			ID:      fmt.Sprintf("%v", node["id"]),
			Name:    fmt.Sprintf("%v", node["name"]),
			Tagline: fmt.Sprintf("%v", node["tagline"]),
		}
		products = append(products, product)
	}
	return products, nil
}

func (ph *ProductHunt) GetProductDetails(slug string) (*Product, error) {
	query := fmt.Sprintf(`
    query {
      post(slug: "%s") {
        name
        tagline
        description
        website
      }
    }
    `, slug)
	data, err := fetchData(query, ph.APIKey)
	if err != nil {
		return nil, err
	}
	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data'")
	}
	post, ok := dataMap["post"].(map[string]interface{})
	if !ok || post == nil {
		// No product found for the given slug.
		return nil, nil
	}
	product := &Product{
		Name:        fmt.Sprintf("%v", post["name"]),
		Tagline:     fmt.Sprintf("%v", post["tagline"]),
		Description: fmt.Sprintf("%v", post["description"]),
		Website:     fmt.Sprintf("%v", post["website"]),
	}
	return product, nil
}

func (ph *ProductHunt) GetPostsByTopic(topic string) ([]Product, error) {
	query := fmt.Sprintf(`
    query {
      posts(order: NEWEST, first: 99, topic: "%s") {
        edges {
          node {
            id
            name
            slug
            tagline
            description
            website
            url
            thumbnail {
              url
            }
          }
        }
      }
    }
    `, topic)
	data, err := fetchData(query, ph.APIKey)
	if err != nil {
		return nil, err
	}
	var products []Product
	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'data'")
	}
	posts, ok := dataMap["posts"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'posts'")
	}
	edges, ok := posts["edges"].([]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'edges'")
	}

	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}

		thumbnailUrl := ""
		if thumb, exists := node["thumbnail"]; exists {
			if thumbMap, ok := thumb.(map[string]interface{}); ok {
				thumbnailUrl = fmt.Sprintf("%v", thumbMap["url"])
			}
		}
		product := Product{
			ID:          fmt.Sprintf("%v", node["id"]),
			Name:        fmt.Sprintf("%v", node["name"]),
			Slug:        fmt.Sprintf("%v", node["slug"]),
			Tagline:     fmt.Sprintf("%v", node["tagline"]),
			Description: fmt.Sprintf("%v", node["description"]),
			Website:     fmt.Sprintf("%v", node["website"]),
			URL:         fmt.Sprintf("%v", node["url"]),
			Thumbnail:   thumbnailUrl,
		}
		products = append(products, product)
	}
	return products, nil
}

func (ph *ProductHunt) GetTopProductsByDate(date string) ([]Product, error) {
	query := fmt.Sprintf(`
    query {
      posts(order: RANKING, postedAfter: "%sT00:00:00Z", postedBefore: "%sT23:59:59Z", first: 5) {
        edges {
          node {
            id
            name
            tagline
            description
            website
          }
        }
      }
    }
    `, date, date)
	data, err := fetchData(query, ph.APIKey)
	if err != nil {
		return nil, err
	}
	var products []Product
	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'data'")
	}
	posts, ok := dataMap["posts"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'posts'")
	}
	edges, ok := posts["edges"].([]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'edges'")
	}
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		product := Product{
			ID:      fmt.Sprintf("%v", node["id"]),
			Name:    fmt.Sprintf("%v", node["name"]),
			Tagline: fmt.Sprintf("%v", node["tagline"]),
		}
		products = append(products, product)
	}
	return products, nil
}

func (ph *ProductHunt) GetProductsByRankByDate(date string, limit int) ([]Product, error) {
	postedAfter := fmt.Sprintf("%sT00:00:00Z", date)
	postedBefore := fmt.Sprintf("%sT23:59:59Z", date)

	query := fmt.Sprintf(`
	query {
	  posts(order: RANKING, postedAfter: "%s", postedBefore: "%s", first: %d) {
	    edges {
	      node {
	        id
	        name
	        tagline
	        description
	        website
			url
	      }
	    }
	  }
	}`, postedAfter, postedBefore, limit)

	data, err := fetchData(query, ph.APIKey)
	if err != nil {
		return nil, err
	}

	var products []Product
	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'data'")
	}

	posts, ok := dataMap["posts"].(map[string]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'posts'")
	}

	edges, ok := posts["edges"].([]interface{})
	if !ok {
		return products, fmt.Errorf("invalid response format: missing 'edges'")
	}

	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]interface{})
		if !ok {
			continue
		}
		node, ok := edgeMap["node"].(map[string]interface{})
		if !ok {
			continue
		}
		product := Product{
			ID:          fmt.Sprintf("%v", node["id"]),
			Name:        fmt.Sprintf("%v", node["name"]),
			Tagline:     fmt.Sprintf("%v", node["tagline"]),
			Description: fmt.Sprintf("%v", node["description"]),
			Website:     fmt.Sprintf("%v", node["website"]),
			URL:         fmt.Sprintf("%v", node["url"]),
		}
		products = append(products, product)
	}

	return products, nil
}
