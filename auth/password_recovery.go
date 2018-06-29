package auth

import (
	"foodcourt/cache"
	"time"
	"math/rand"
	"strconv"
	"foodcourt/database"
	"fmt"
	"foodcourt/config"
	"foodcourt/logger"
	"foodcourt/mails"
	"foodcourt/environment"
)

// Password recovery code TTL in seconds
const recoveryCodeSecondsTTL = 60
const recoveryCodeLength = 6
const restoreEmailFileName = "restore-password-mail.html"

type RecoveryCodeMail struct {
	Code string
}

func getResetCodeKey(email string) string {
	return "resetcode__" + EncryptString(email)
}

// ResetPassword resets user password by core token retrieved from previous recovery step and client info
func ResetPassword(newPassword string, coreToken string, userAgent string, ip string) error {
	tokenPair := NewTokenPair(coreToken, userAgent, ip)

	exists, err := tokenPair.Exists()

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Bad token")
	}

	email, err := tokenPair.GetTokenOwner()
	defer tokenPair.Reveal()

	if err != nil {
		return err
	}

	db := database.GetInstance()
	defer db.Close()

	err, userId := GetIdForEmail(db, email)

	if err != nil {
		return err
	}

	user := NewUser()
	user.ID = userId
	user.Password = newPassword

	return UpdateUser(db, &user)
}

// Creates and registers reset token pair
func CreateResetTokenFromCode(email string, code string, agent string, ip string) (string, error) {
	cache.Client.Del(getResetCodeKey(email)).Result()
	tokenPair := BuildResetTokenPair(email, code, agent, ip)

	err := tokenPair.Save(email)

	return tokenPair.CoreKey, err
}

// CanRequestCode checks if the user can request a new code
func CanRequestCode(email string) (bool, string, error) {
	db := database.GetInstance()

	if err, userExists := MailExists(db, email); err != nil {
		return false, "", err
	} else if !userExists {
		return false, "No such user with email " + email, nil
	}

	defer db.Close()

	if hasCode, err := HasResetCode(email); err != nil {
		return false, "", err
	} else if hasCode {
		return false, "You've already requested a code, check your email or try again later", nil
	}

	return true, "", nil
}

func generateRecoveryCode() (code string) {
	count := 1
	for {
		code += strconv.Itoa(rand.Intn(9))
		if count == recoveryCodeLength { break }
		count++
	}
	return code
}

func SendRestoreCode(email string) error {
	code, err := createResetCode(email)

	if err != nil {
		return err
	}

	// Just print code to log if debug mode is enabled
	if config.GetBool("DEBUG") {
		logger.GetLogger().Debug("Reset code is %s", code)
		return nil
	}

	mailSender, err := mails.NewMailSender(
		environment.GetResourcePath(restoreEmailFileName),
	);

	if err != nil {
		return err
	}

	mailSender.SetSubject("FoodCount - Your password reset code")

	if err = mailSender.Init(); err != nil {
		return err
	}

	return mailSender.Send(RecoveryCodeMail{code}, mails.NewRecipient(email, "User"));
}

func ResetCodeValid(email string, code string) (bool, error) {
	key := getResetCodeKey(email)
	correctCode, err := cache.Client.Get(key).Result()

	if err != nil {
		return false, err
	}

	isValid := code == correctCode
	return isValid, nil
}

// HasResetCode checks if password reset code was generated for email
func HasResetCode(email string) (bool, error) {
	key := getResetCodeKey(email)
	val, err := cache.Client.Exists(key).Result()

	if err != nil {
		return false, err
	}

	return (val > 0), nil
}

// CreateResetCode generates a new password reset code for email
func createResetCode(email string) (code string, err error) {
	key := getResetCodeKey(email)
	ttl := time.Duration(time.Second * recoveryCodeSecondsTTL)
	code = generateRecoveryCode()

	_, err = cache.Client.Set(key, code, ttl).Result()

	return code, err
}