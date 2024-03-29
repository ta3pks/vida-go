package models

import (
	"fmt"
	"time"

	"gitlab.mugsoft.io/vida/go-api/helpers"
	"gopkg.in/mgo.v2/bson"
)

const _COL_USER_STR = "users"

var _col_user = db_get().C(_COL_USER_STR)

type User struct {
	//{{{
	Id            string    `bson:"id" json:"id"`
	Name          string    `bson:"name" json:"name"`
	Lastname      string    `bson:"lastname" json:"lastname"`
	Phone         string    `bson:"phone" json:"phone"`
	Email         string    `bson:"email" json:"email"`
	Notification  int       `bson:"notification" json:"notification"`
	FbAccountName string    `json:"fb_account_name" bson:"fb_account_name"`
	FbProfilePic  string    `json:"fb_profile_pic" bson:"fb_profile_pic"`
	Password      string    `bson:"password" json:"-"`
	Login_expires time.Time `bson:"-" json:"-"`
	Token         string    `bson:"-" json:"token"`
	ProfilePicURL string    `bson:"profile_pic_url" json:"profile_pic_url"`
	PassReset     bool      `bson:"pass_reset" json:"pass_reset,omitempty"`
	Tmp           bool      `json:"tmp,omitempty" bson:"tmp"`
	Defaults
	//}}}
}

//User_new generates id and date fields of the user and hashes password then saves
func User_new(u *User) error {
	//{{{
	//{{{ error checks
	if nil == u {
		return fmt.Errorf("user cannot be empty")
	}
	if "" == u.Email && "" == u.Phone {
		return fmt.Errorf("missing email and phone")
	}
	var usr = new(User)
	usr.Email = u.Email
	usr.Phone = u.Phone
	if nil == User_get(usr) && !usr.Tmp {
		return fmt.Errorf("user exists")
	} //}}}
	u.Id = helpers.Unique_id()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if usr.Tmp && usr.Id != "" {
		usr.Password = u.Password
		usr.Password = Hash_password(usr, usr.Password)
		usr.UpdatedAt = u.UpdatedAt
		usr.Tmp = false
		helpers.Log(helpers.INFO, "tmp user registering : ", usr.Email)
		err := _col_user.Update(bson.M{"id": usr.Id}, map[string]interface{}{
			"$set": map[string]interface{}{
				"name":         u.Name,
				"lastname":     u.Lastname,
				"phone":        u.Phone,
				"email":        u.Email,
				"notification": u.Notification,
				"password":     usr.Password,
				"tmp":          false,
			},
		})
		return err
	}
	u.Password = Hash_password(u, u.Password)
	return _col_user.Insert(u)
	//}}}
}

//Hash_password hashes user password with salt and generates id if it is empty
func Hash_password(u *User, pass string) string {
	//{{{
	if "" == u.Id {
		u.Id = helpers.Unique_id()
	}
	return helpers.MD5(u.Id + pass)
	//}}}
}
func User_get_by_id(id string) (*User, error) {
	//{{{
	var usr = new(User)
	err := _col_user.Find(map[string]string{"id": id}).One(usr)
	return usr, err
	//}}}
}

func User_get(u *User) error {
	//{{{
	_q := []bson.M{}
	if "" != u.Email {
		_q = append(_q, bson.M{"email": u.Email})
	}
	if "" != u.Phone {
		_q = append(_q, bson.M{"phone": u.Phone})
	}
	if "" == u.Email && "" == u.Phone {
		return fmt.Errorf("missing email and phone")
	}
	return _col_user.Find(bson.M{
		"$or": _q,
	}).One(u)
	//}}}
}
func User_get_by_email(email string) (*User, error) {
	//{{{
	if !helpers.Is_email_valid(email) {
		return nil, fmt.Errorf("invalid email address")
	}
	usr := &User{
		Email: email,
	}
	err := User_get(usr)
	if nil != err {
		return nil, err
	}
	return usr, nil
	//}}}
}
func User_update(userid string, fields map[string]interface{}, updatedU *User) error {
	//{{{
	var _fields_with_pdatedAt = map[string]interface{}{
		"updated_at": time.Now(),
	}
	for k, v := range fields {
		_fields_with_pdatedAt[k] = v
	}
	err := _col_user.Update(bson.M{"id": userid}, bson.M{"$set": _fields_with_pdatedAt})
	if nil != err {
		return err
	}
	if nil == updatedU {
		return nil
	}
	updatedU.Id = userid
	return User_get(updatedU)
	//}}}
}
func User_new_tmp(email string) (*User, error) {
	//{{{
	if !helpers.Is_email_valid(email) {
		return nil, fmt.Errorf("invalid email address")
	}
	u := &User{
		Email: email,
	}
	//{{{ error checks
	if nil == User_get(u) {
		return nil, fmt.Errorf("user exists")
	} //}}}
	u.Tmp = true
	err := User_new(u)
	if nil != u {
		u.Token = helpers.Unique_id()
	}
	return u, err //}}}
}

//User_or_tmp returns the given user by email or creates and returns a tmp user
func User_or_tmp(email string) (*User, error) {
	//{{{
	if !helpers.Is_email_valid(email) {
		return nil, fmt.Errorf("invalid email address")
	}
	u := &User{
		Email: email,
	}
	if nil == User_get(u) {
		u.Token = helpers.Unique_id()
		return u, nil
	}
	u.Tmp = true
	err := User_new(u)
	if nil != u {
		u.Token = helpers.Unique_id()
	}
	return u, err //}}}
}
func User_delete(email string) error {
	//{{{
	return _col_user.Remove(bson.M{"email": email}) //}}}
}
