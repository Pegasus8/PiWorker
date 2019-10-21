package configs

import (
)

// NewUser is the function used to add a new user.
func NewUser(username, password string, admin bool) error {
	if usernameExists(username) {
		return ErrUsernameExists
	}
	newUser := User{
		username, 
		password,
		admin,
	}
	CurrentConfigs.Users = append(CurrentConfigs.Users, newUser)
	err := WriteToFile()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser is used to delete a existing users.
func DeleteUser(username string) error {
	// TODO Verify if the user that executes this function (from the API)
	// is an admin user.
	for index, user := range CurrentConfigs.Users {
		if user.Username == username {
			CurrentConfigs.Users = append(
				CurrentConfigs.Users[:index], CurrentConfigs.Users[index+1:]...
			)
			err := WriteToFile()
			if err != nil {
				return err
			}
			
			return nil
		}
	}

	return ErrUserNotFound
}

// AuthUser is used to authenticate a user.
func AuthUser(username, password string) (authenticated bool) {
	for _, user := range CurrentConfigs.Users {
		// TODO Hashed & encrypted password 
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return false
			}
			return true
		}
	}
	return false
}

// ChangeUserPassword is a function used to change the password of a user.
func ChangeUserPassword(username, newPassword string) error {
	for _, user := range CurrentConfigs.Users {
		if user.Username == username {
			// TODO Hashed & encrypted password
			hashedPwd, err := hashAndSalt([]byte(newPassword))
			if err != nil {
				return err
			}
			user.Password = hashedPwd
			err = WriteToFile()
			if err != nil {
				return err
			}
			return nil
		}
	}

	return ErrUserNotFound
}

func usernameExists(username string) bool {
	for _, user := range CurrentConfigs.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func hashAndSalt(password []byte) (hashedPassword string, err error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}