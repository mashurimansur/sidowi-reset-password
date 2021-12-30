package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mashurimansur/sidowi-reset-password/database"
	"golang.org/x/crypto/bcrypt"
)

type Kaders struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	Name           string     `gorm:"column:name;type:varchar(100);not null;" json:"name"`
	Email          string     `gorm:"column:email;type:varchar(255);unique;not null;" json:"email"`
	NIK            *string    `gorm:"column:nik;type:varchar(16);unique" json:"nik"`
	DateBirth      time.Time  `gorm:"column:date_birth;type:date;not null;" json:"date_birth"`
	PlaceBirth     string     `gorm:"column:place_birth;type:varchar(50);not null;" json:"place_birth"`
	Avatar         string     `gorm:"column:avatar;type:varchar(255);not null;" json:"avatar"`
	Job            string     `gorm:"column:job;type:varchar(100);not null;" json:"job"`
	Office         string     `gorm:"column:office;type:varchar(100);not null;" json:"office"`
	Skills         string     `gorm:"column:skills;type:varchar(255);not null;" json:"skills"`
	Address        string     `gorm:"column:address;type:varchar(255);not null;" json:"address"`
	Phone          string     `gorm:"column:phone;type:varchar(15);unique;not null;" json:"phone"`
	BloodType      string     `gorm:"column:blood_type;type:varchar(2);not null;" json:"blood_type"`
	Gender         string     `gorm:"column:gender;type:varchar(1);not null;" json:"gender"`
	ZipCode        string     `gorm:"column:zip_code;type:varchar(7);not null;" json:"zip_code"`
	ProvinceID     string     `gorm:"column:province_id;type:char(2);not null;" json:"province_id"`
	CityID         string     `gorm:"column:city_id;type:char(4);not null;" json:"city_id"`
	DistrictID     string     `gorm:"column:district_id;type:char(7);not null;" json:"district_id"`
	VillageID      string     `gorm:"column:village_id;type:char(10);not null;" json:"village_id"`
	RegistrationID *uint      `gorm:"column:registration_id;" json:"registration_id"`
	Password       string     `gorm:"column:password;type:varchar(255);not null;" json:"password"`
	Status         string     `gorm:"column:status;type:varchar(30)" json:"status"`
	CampusName     string     `gorm:"column:campus_name;type:varchar(30)" json:"campus_name"`
	CampusMajor    string     `gorm:"column:campus_major;type:varchar(30)" json:"campus_major"`
	CampusBatch    string     `gorm:"column:campus_batch;type:varchar(4)" json:"campus_batch"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `sql:"index" json:"deleted_at"`
}

var (
	postgres *gorm.DB
)

func init() {
	database.LoadEnv()
	postgres = database.ConnectPostgres()

}

func main() {
	kaders, err := findKader()
	if err != nil {
		fmt.Println(err.Error())
	}

	bar := pb.StartNew(len(kaders))
	for _, v := range kaders {
		password := generatePassword(v.DateBirth.Format("02012006"))
		v.Password = password
		v.UpdatedAt = time.Now()
		errUpdate := updateKader(&v)
		if errUpdate != nil {
			fmt.Println(fmt.Printf("%d %s Gagal update password", v.ID, v.Name))
			continue
		} else {
			bar.Increment()
			time.Sleep(time.Millisecond)
		}
	}
	bar.Finish()
	fmt.Println("Successfully set new password")
}

func findKader() (kaders []Kaders, err error) {
	if err := postgres.Table("kaders").Where("email != ?", "mashurimansur18@gmail.com").Find(&kaders).Error; err != nil {
		return kaders, err
	}
	return
}

func updateKader(kader *Kaders) (err error) {
	if err := postgres.Table("kaders").Model(kader).Updates(kader).Error; err != nil {
		return err
	}
	return
}

func generatePassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
