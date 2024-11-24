package filter

import "github.com/VsProger/snippetbox/internal/models"

func (f *FilterRepo) getAllCategoriesByPostId(id int) ([]models.Category, error) {
	query := `
		SELECT c.ID, c.Name
		FROM Category c
		INNER JOIN PostCategory pc ON c.ID = pc.CategoryID
		WHERE pc.PostID = ?
	`
	rows, err := f.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
