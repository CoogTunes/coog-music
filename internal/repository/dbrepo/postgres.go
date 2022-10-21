package dbrepo

func (n *postgresDBRepo) AddUser(res models.user) err {
	query := "insert into users (username, password) values ($1, $2)"

	_, err := m.DB.Exec(query, rest.username, res.password)
	if err != nil {
		return err
	}

	return nil
}
