package user

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"strings"
	"errors"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	pb "gokitgrpc/user/pb"
)

type (
	CreateUserRequest struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	CreateUserResponse struct {
		Ok string `json:"ok"`
	}

	UpdateUserRequest struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}
	UpdateUserResponse struct {
		Ok string `json:"ok"`
	}

	DeleteUserRequest struct {
		Id       string `json:"id"`
	}
	DeleteUserResponse struct {
		Ok string `json:"ok"`
	}

	ViewUserRequest struct {
		Id string `json:"id"`
	}
	ViewUserResponse struct {
		User    string `json:"user"`
	}

	ListUserRequest struct {
		Limit string `json:"limit"`
		Offset string `json:"offset"`
	}

	AuthenticateUserRequest struct {
		Phone       string `json:"phone"`
		Password    string `json:"password"`
	}

	AuthenticateUserResponse struct {
		AccessToken string `json:"access_token"`
	}
)

type AccessDetails struct {
    UserId   string
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
	   return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
	   return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
	   return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
	   return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)

	if err != nil {
	   return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
	   user_id, ok := claims["user_id"].(string)
	   if !ok {
		  return nil, err
	   }
	   
	   return &AccessDetails{
		  UserId:   user_id,
	   }, nil
	}

	return nil, err
}

func secureEndpoint(r *http.Request) error {
	auth, err := ExtractTokenMetadata(r)
	if err != nil {
		return errors.New("Token not valid")//+err.Error())
	}
	if (auth != nil) {
 		fmt.Println("Login as : "+auth.UserId)
	} else {
		return errors.New("User not authenticated")
	}

	return errors.New("Not Authorized")
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeViewUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ViewUserRequest
	vars := mux.Vars(r)

	req = ViewUserRequest{
		Id: vars["id"],
	}
	return req, nil
}

func decodeListUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ListUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeDeleteUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	err := secureEndpoint(r)
	if (err != nil) {
		return nil, err
	}
	var req DeleteUserRequest
	vars := mux.Vars(r)

	req = DeleteUserRequest{
		Id: vars["id"],
	}
	return req, nil
}

func decodeUpdateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	err := secureEndpoint(r)
	if (err != nil) {
		return nil, err
	}
	var req UpdateUserRequest
	vars := mux.Vars(r)
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	req.Id = vars["id"]
	return req, nil
}

func decodeAuthenticateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req AuthenticateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}


func decodeGViewUserReq(_ context.Context, request interface{}) (interface{}, error) {
    grpc := request.(*pb.ViewUserRequest)
	return grpc, nil
}

func encodeGUser(_ context.Context, response interface{}) (interface{}, error) {
    //resp := response.(pb.User)
    return &pb.UserResponse{}, nil
}

func decodeGCreateUserReq(_ context.Context, request interface{}) (interface{}, error) {
    grpc := request.(*pb.CreateUserRequest)
	return grpc, nil
}

func encodeGResponse(_ context.Context, response interface{}) (interface{}, error) {
    //resp := response.(pb.User)
    return &pb.Response{}, nil
}