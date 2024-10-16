package models

type Category struct {
	ID   int
	Name string
}

//
//func (m *PostModel) GetCategoryIDs() ([]int, error) {
//	stmt := `SELECT ID FROM Category`
//	rows, err := m.DB.Query(stmt)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var categoryIDs []int
//	for rows.Next() {
//		var cID int
//		if err := rows.Scan(&cID.ID); err != nil {
//			return nil, err
//		}
//		categoryIDs = append(categoryIDs, cID)
//	}
//	return categoryIDs, nil
//}
