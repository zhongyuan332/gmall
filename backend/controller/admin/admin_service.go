package admin

import (
	"database/sql"
	"errors"
	"github.com/zhongyuan332/gmall/backend/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	DB *sql.DB
}

// GetByID 通过ID获取用户
func (s *UserService) GetByID(id int64) (*model.AdminUser, error) {
	user := &model.AdminUser{}

	query := `SELECT id, username, password, email, real_name, mobile, avatar, 
		role, status, last_login_time, last_login_ip, login_count, created_at, updated_at 
		FROM admin_user WHERE id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.RealName,
		&user.Mobile, &user.Avatar, &user.Role, &user.Status, &user.LastLoginTime,
		&user.LastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

// GetByUsername 通过用户名获取用户
func (s *UserService) GetByUsername(username string) (*model.AdminUser, error) {
	user := &model.AdminUser{}

	query := `SELECT id, username, password, email, real_name, mobile, avatar, 
		role, status, last_login_time, last_login_ip, login_count, created_at, updated_at 
		FROM admin_user WHERE username = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.RealName,
		&user.Mobile, &user.Avatar, &user.Role, &user.Status, &user.LastLoginTime,
		&user.LastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

// GetAll 获取所有用户
func (s *UserService) GetAll() ([]*model.AdminUser, error) {
	query := `SELECT id, username, email, real_name, mobile, avatar, 
		role, status, last_login_time, last_login_ip, login_count, created_at, updated_at 
		FROM admin_user WHERE deleted_at IS NULL ORDER BY id DESC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.AdminUser
	for rows.Next() {
		user := &model.AdminUser{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.RealName,
			&user.Mobile, &user.Avatar, &user.Role, &user.Status, &user.LastLoginTime,
			&user.LastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *model.AdminUser) error {
	// 检查用户是否存在
	_, err := s.GetByID(user.ID)
	if err != nil {
		return err
	}

	var query string
	var args []interface{}

	// 如果提供了新密码，则需要加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		query = `UPDATE admin_user SET 
			username = ?, password = ?, email = ?, real_name = ?, 
			mobile = ?, role = ?, status = ? 
			WHERE id = ? AND deleted_at IS NULL`

		args = []interface{}{
			user.Username, string(hashedPassword), user.Email, user.RealName,
			user.Mobile, user.Role, user.Status, user.ID,
		}
	} else {
		// 没有提供新密码，保持原密码不变
		query = `UPDATE admin_user SET 
			username = ?, email = ?, real_name = ?, 
			mobile = ?, role = ?, status = ? 
			WHERE id = ? AND deleted_at IS NULL`

		args = []interface{}{
			user.Username, user.Email, user.RealName,
			user.Mobile, user.Role, user.Status, user.ID,
		}
	}

	_, err = s.DB.Exec(query, args...)
	return err
}

// DeleteUser 软删除用户
func (s *UserService) DeleteUser(id int64) error {
	query := `UPDATE admin_user SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := s.DB.Exec(query, time.Now(), id)
	return err
}

// UpdateLoginInfo 更新用户登录信息
func (s *UserService) UpdateLoginInfo(id int64, ip string) error {
	query := `UPDATE admin_user SET 
		last_login_time = ?, last_login_ip = ?, login_count = login_count + 1 
		WHERE id = ?`

	_, err := s.DB.Exec(query, time.Now(), ip, id)
	return err
}

// VerifyPassword 验证用户密码
func (s *UserService) VerifyPassword(username, password string) (*model.AdminUser, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if !user.Status {
		return nil, errors.New("用户已禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(user *model.AdminUser) error {
	// 检查用户名是否已存在
	_, err := s.GetByUsername(user.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 插入新用户
	query := `INSERT INTO admin_user (username, password, email, real_name, mobile, role, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := s.DB.Exec(
		query,
		user.Username,
		string(hashedPassword),
		user.Email,
		user.RealName,
		user.Mobile,
		user.Role,
		user.Status,
	)

	if err != nil {
		return err
	}

	// 获取生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}
