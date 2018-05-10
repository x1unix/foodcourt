package controller

import (
	"net/http"
	"io"
	"foodcourt/auth"
	"encoding/json"
	"foodcourt/rest"
	"fmt"
)

func PasswordRecovery(w http.ResponseWriter, r *http.Request) {
	req, err := readRecoveryAction(&r.Body)

	if err != nil {
		rest.BadRequest(&w, "Failed to read request payload: " + err.Error())
		return
	}

	switch req.Action {
	case auth.ActionRequestCode:
		handleCodeRequest(&w, req)
		return
	case auth.ActionSubmitCode:
		handleCodeSubmit(r, &w, req)
		return
	case auth.ActionSetPassword:
		handlePasswordReset(r, &w, req)
	default:
		rest.BadRequest(&w, fmt.Sprintf("Unknown action '%s'", req.Action))
		return
	}
}

func handlePasswordReset(r *http.Request, w *http.ResponseWriter, action *auth.RecoveryAction) {
	var coreToken string
	var password string
	var ok bool

	if coreToken, ok = action.Data["resetToken"]; !ok {
		rest.BadRequest(w, "Bad Payload")
		return
	}

	if password, ok = action.Data["newPassword"]; !ok {
		rest.BadRequest(w, "Bad Payload")
		return
	}

	if err := auth.ResetPassword(password, coreToken, r.Header.Get("User-Agent"), rest.GetClientIP(r)); err != nil {
		log.Error("failed to reset password for token %s: %v", coreToken, err)
		rest.ErrorFromString("Cannot reset password, internal error", 500).Write(w)
	}
}

func handleCodeSubmit(r *http.Request, w *http.ResponseWriter, action *auth.RecoveryAction) {
	var email string
	var code string
	var ok bool

	if email, ok = action.Data["email"]; !ok {
		rest.BadRequest(w, "Email required")
		return
	}

	if code, ok = action.Data["code"]; !ok {
		rest.BadRequest(w, "No code provided")
		return
	}

	valid, err := auth.ResetCodeValid(email, code);

	if err != nil {
		rest.Error(fmt.Errorf("Cannot check if reset code is valid, internal error")).Write(w)
		log.Error("cannot check if reset code of %s is valid: %v", email, err)
		return
	}

	if !valid {
		rest.BadRequest(w, "Invalid or expired code")
		return
	}

	token, err := auth.CreateResetTokenFromCode(
		email,
		code,
		r.Header.Get("User-Agent"),
		rest.GetClientIP(r),
	)

	if err != nil {
		log.Error("cannot create reset password token for user %s by code %s: %v", email, code, err)
		rest.Error(fmt.Errorf("internal error, please try again later")).Write(w)
		return
	}

	rest.Success(token).Write(w)
}

func handleCodeRequest(w *http.ResponseWriter, action *auth.RecoveryAction) {
	if email, ok := action.Data["email"]; !ok {
		rest.BadRequest(w, "No email provided")
		return
	} else {
		can, message, err := auth.CanRequestCode(email)

		if err != nil {
			rest.Error(err).Write(w)
			return
		}

		if !can {
			rest.BadRequest(w, message)
			return
		}

		if err := auth.SendRestoreCode(email); err != nil {
			log.Error("failed to send recovery code to '%s': %v", email, err)
			rest.ErrorFromString("Failed to send recovery email, internal error", 500).Write(w)
			return
		}

		rest.Ok(w)
		return
	}
}

func readRecoveryAction(bodyPtr *io.ReadCloser) (*auth.RecoveryAction, error) {
	var out auth.RecoveryAction
	body := *bodyPtr
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&out)
	defer body.Close()
	return &out, err
}
