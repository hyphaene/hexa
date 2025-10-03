package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// UpdateYAMLField met à jour un champ YAML en préservant commentaires et structure
// Supporte les chemins imbriqués avec notation pointée: "jira.boardId"
func UpdateYAMLField(filePath string, key string, value any) error {
	// Lire le fichier existant (ou créer structure vide)
	var root yaml.Node
	data, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("reading file: %w", err)
		}
		// Fichier n'existe pas, créer un document vide
		root.Kind = yaml.DocumentNode
		root.Content = []*yaml.Node{{
			Kind: yaml.MappingNode,
		}}
	} else if len(data) == 0 {
		// Fichier vide, créer un document vide
		root.Kind = yaml.DocumentNode
		root.Content = []*yaml.Node{{
			Kind: yaml.MappingNode,
		}}
	} else {
		if err := yaml.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parsing yaml: %w", err)
		}
	}

	// Naviguer et mettre à jour le champ
	if err := setNestedField(&root, key, value); err != nil {
		return fmt.Errorf("setting field: %w", err)
	}

	// Écrire avec préservation
	output, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("marshaling yaml: %w", err)
	}

	return os.WriteFile(filePath, output, 0644)
}

// setNestedField gère la notation pointée "parent.child" et crée les nœuds manquants
func setNestedField(root *yaml.Node, key string, value any) error {
	// Trouver le mapping root (premier élément de DocumentNode)
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return fmt.Errorf("invalid yaml structure")
	}

	mappingNode := root.Content[0]
	if mappingNode.Kind != yaml.MappingNode {
		return fmt.Errorf("root is not a mapping")
	}

	// Support notation pointée: "jira.boardId"
	keys := splitKey(key)
	return setNestedFieldRecursive(mappingNode, keys, value)
}

// splitKey sépare "jira.boardId" en ["jira", "boardId"]
func splitKey(key string) []string {
	var result []string
	current := ""
	for _, char := range key {
		if char == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// setNestedFieldRecursive navigue récursivement dans les clés
func setNestedFieldRecursive(mappingNode *yaml.Node, keys []string, value any) error {
	if len(keys) == 0 {
		return fmt.Errorf("empty key path")
	}

	currentKey := keys[0]
	isLastKey := len(keys) == 1

	// Chercher la clé dans le mapping actuel
	for i := 0; i < len(mappingNode.Content); i += 2 {
		keyNode := mappingNode.Content[i]
		valueNode := mappingNode.Content[i+1]

		if keyNode.Value == currentKey {
			if isLastKey {
				// Dernière clé: remplacer la valeur
				newValueNode := &yaml.Node{}

				if err := newValueNode.Encode(value); err != nil {
					return fmt.Errorf("encoding value: %w", err)
				}
				mappingNode.Content[i+1] = newValueNode
				return nil
			} else {
				// Clé intermédiaire: descendre dans le mapping
				if valueNode.Kind != yaml.MappingNode {
					// Transformer en mapping si ce n'en est pas un
					valueNode.Kind = yaml.MappingNode
					valueNode.Content = []*yaml.Node{}
				}
				return setNestedFieldRecursive(valueNode, keys[1:], value)
			}
		}
	}

	// Clé non trouvée: créer
	if isLastKey {
		// Dernière clé: ajouter directement
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: currentKey,
		}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(value); err != nil {
			return fmt.Errorf("encoding value: %w", err)
		}
		mappingNode.Content = append(mappingNode.Content, keyNode, valueNode)
		return nil
	} else {
		// Clé intermédiaire: créer mapping et descendre
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: currentKey,
		}
		newMappingNode := &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: []*yaml.Node{},
		}
		mappingNode.Content = append(mappingNode.Content, keyNode, newMappingNode)
		return setNestedFieldRecursive(newMappingNode, keys[1:], value)
	}
}

// ReadYAMLField lit un champ depuis un fichier YAML
func ReadYAMLField(filePath string, key string) (any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	value, ok := config[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return value, nil
}
