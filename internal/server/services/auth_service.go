package services

type AuthService struct {
	UserRepo repositories.UserRepository
}

// Verifica la cedula y contraseña del usuario
func (as *AuthService) LoginUser(cedula string, password string) (string, error) {

	// Busca el usuario por cedula en el repositorio
	user, err := as.UserRepo.FindByCedula(cedula)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	// Verifica que la contraseña coincida (inicialmente igual a la cédula)
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	// Genera un token (lógica que deberías implementar para JWT u otro tipo de token)
	token := "generated-token"  // Reemplazar con tu lógica de generación de tokens

	return token, nil
}

// Cambia la contraseña de un usuario
func (as *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := as.UserRepo.FindByID(userID)
	if err != nil || user == nil {
			return errors.New("user not found")
	}

	// Verifica que la contraseña antigua coincida
	if user.Password != oldPassword {
			return errors.New("incorrect old password")
	}

	// Actualiza la contraseña
	user.Password = newPassword
	return as.UserRepo.UpdateUser(user)  
}
