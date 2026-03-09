package database

import (
	"ADS4/internal/models"
	_ "database/sql"
)

// GetAllUsers function
func (db *DB) GetAllUsers() ([]models.User, error) {
	query := `SELECT userid, username, email, role, defaultadmin, active FROM userT`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.DefaultAdmin,
			&user.Active,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Create user function
func (db *DB) CreateUser(user *models.User) error {
	query := `
		INSERT INTO userT (username, password, email, role)
		VALUES ($1, $2, $3, $4)
		`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(user.Username, user.Password, user.Email, user.Role)

	if err != nil {
		return err
	}

	return nil
}

// Update user function
func (db *DB) UpdateUserWithPassword(user *models.User) error {
	query := `
        UPDATE userT
        SET username = $1, email = $2, role = $3, password = $4, active = $5
        WHERE userid = $6
        `
	args := []interface{}{user.Username, user.Email, user.Role, user.Password, user.Active, user.UserID}

	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

// Update user function
func (db *DB) UpdateUser(user *models.User) error {
	query := `
		UPDATE userT
		SET username = $1, email = $2, role = $3, active = $4
		WHERE userid = $5
		`

	args := []interface{}{user.Username, user.Email, user.Role, user.Active, user.UserID}

	updateStmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(args...)

	if err != nil {
		return err
	}

	return nil

}

// Get user by username function
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role, defaultadmin, active
		FROM userT
		WHERE username = $1
		`
	var user models.User
	err := db.QueryRow(query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.DefaultAdmin,
		&user.Active,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by ID function
func (db *DB) GetUserByID(userid int) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role, defaultadmin, active
		FROM userT
		WHERE userid = $1
		`

	var user models.User
	err := db.QueryRow(query).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.DefaultAdmin,
		&user.Active,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete user function
func (db *DB) DeleteUser(userid int) error {
	query := `DELETE FROM userT WHERE userid = $1`
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(userid)

	if err != nil {
		return err
	}

	return nil
}

// Update password function
func (db *DB) UpdatePassword(userid int, password string) error {
	query := `
		UPDATE userT
		SET password = $1
		WHERE userid = $2
		`
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(password, userid)

	if err != nil {
		return err
	}

	return nil
}

// Update password function
func (db *DB) UpdateActive(userid int, active bool) error {
	query := `
		UPDATE userT
		SET active = $1
		WHERE userid = $2
		`
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(active, userid)

	if err != nil {
		return err
	}

	return nil
}

// Get user by email function
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role,active
		FROM userT
		WHERE email = $1
		`
	var user models.User
	err := db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.Active,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// check if the suer is active
func (db *DB) IsUserActive(userid int) bool {
	query := `
		SELECT active
		FROM userT
		WHERE userid = $1 AND active = 1
		`
	var user models.User
	err := db.QueryRow(query, userid).Scan(
		&user.Active,
	)

	if err != nil {
		return false
	}

	return true
}
