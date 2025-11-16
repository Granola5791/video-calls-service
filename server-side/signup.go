package main

func SignUpUser(username string, password string) error {
	hashedPassword, salt, err := GenerateNewHashedPassword(password)
	if err != nil {
		return err
	}
	err = InsertUserToDB(username, hashedPassword, salt)
	if err != nil {
		return err
	}

	return nil
}