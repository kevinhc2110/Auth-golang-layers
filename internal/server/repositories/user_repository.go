package repositories

// Define los metodos para manejar la autenticacion de los usuarios
type AuthRepository  interface {
	FindUserByID(userID string) (*models.User, error)
	ChangePassword(userID, oldPassword, newPassword string) error
}

//Implementacion concreta de AuthRepository
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository crea una nueva instancia de AuthRepository
func NewUserRepository(db *sql.DB) AuthRepository  {
	return &authRepository{db: db}
}

// FindUserByID busca un usuario por su ID
func (r *authRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, password FROM users WHERE id = $1`
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
		return nil, nil
		}
	return nil, err
	}
	return &user, nil
}

// ChangePassword cambia la contraseña del usuario
func (r *authRepository) ChangePassword(userID, oldPassword, newPassword string) error {
	// Primero, verifica si la contraseña antigua es correcta
	user, err := r.FindUserByID(userID)
	if err != nil {
			return err
	}
	if user == nil {
		return errors.New("Usuario no encontrado") 
	}

	// Actualiza la contraseña del usuario
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("Contraseña antigua incorrecta") 

	 // Cifra la nueva contraseña
	 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	 if err != nil {
			 return err 
	 }

	 // Actualiza la contraseña del usuario
	 query := "UPDATE users SET password = $1 WHERE cedula = $2"
	 _, err = r.db.Exec(query, hashedPassword, userID)
	 return err

}