package repository

import (
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/domain"
	"database/sql"
	"fmt"
)

// CreateUser はハッシュ化されたパスワードで新しいユーザーをデータベースに保存します。
func CreateUser(username string, password string, email string) (int64, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("repository.CreateUser: %w", err)
	}

	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"

	result, err := database.DB.Exec(query, username, hashedPassword, email)
	if err != nil {
		return 0, fmt.Errorf("could not insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not retrieve last insert ID: %v", err)
	}

	return id, nil
}

// GetUserByID
func GetUserByID(id int64) (*domain.User, error) {
	query := "SELECT id, username, password, email, created_at, updated_at FROM users WHERE id = ?"

	row := database.DB.QueryRow(query, id)

	var u domain.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByID: could not retrieve user with id %d: %v", id, err)
	}
	return &u, nil
}

// GetAllUsers
func GetAllUsers() ([]domain.User, error) {
	query := "SELECT id, username, password, email, created_at, updated_at FROM users"

	// db.Query
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: could not retrieve users: %v", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetAllUsers: error scanning user row: %v", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: error iterating user rows: %v", err)
	}

	return users, nil
}

//名前でユーザーを取得する
func GetUserByUsername(username string) (*domain.User, error) {
	query := "SELECT id, username, password, email, created_at, updated_at,deleted_at FROM users WHERE username= ? AND deleted_at IS NULL"

	row := database.DB.QueryRow(query, username)

	var u domain.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository.GetUserByUsername: データベースクエリエラー: %w", err)
	}

	return &u, nil

}

func UpdateUserEmail(id int64, newEmail string) (int64, error) {
	query := "UPDATE users SET email = ? WHERE id = ?"
	result, err := database.DB.Exec(query, newEmail, id)
	if err != nil {
		return 0, fmt.Errorf("UpdateUserEmail: could not update user email for id %d: %v", id, err)
	}

	// 获取受该 UPDATE 语句影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("UpdateUserEmail: could not get rows affected after update: %v", err)
	}

	// 返回受影响的行数和 nil 错误
	return rowsAffected, nil
}

// 更新置顶用户密码
func UpdateUserPassword(newPassword string, id int64) (int64, error) {
	query := "UPDATE users SET password = ?  WHERE id = ?"
	result, err := database.DB.Exec(query, newPassword, id)
	if err != nil {
		return 0, fmt.Errorf("UpdateUserPassword: could not update user password for id %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("UpdateUserPassword: could not get rows affected after update: %v", err)
	}
	return rowsAffected, nil
}

// DeleteUser 根据 ID 从数据库中删除用户
func DeleteUser(id int64) (int64, error) {
	query := "DELETE FROM users WHERE id = ?"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("DeleteUser: could not delete user with id %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("DeleteUser: could not get rows affected after delete: %v", err)
	}

	// 返回受影响的行数和 nil 错误
	return rowsAffected, nil
}
