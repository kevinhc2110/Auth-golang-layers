package models

// LoginRequest - Representa las credenciales del usuario para el login
type LoginRequest struct {
    Id   string `json:"id" binding:"required"`  
    Password string `json:"password" binding:"required"` 
}

// ChangePasswordRequest - Representa la solicitud de cambio de contraseña
type ChangePasswordRequest struct {
    UserID      uint   `json:"user_id" binding:"required"`
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required"`
}



