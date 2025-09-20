package validator

import (
	"reflect"
	"strings"
)

func registerTagsCustomFields() {
	validate.RegisterTagNameFunc(persianTagNamesEvaluator)
}

// HELPERS

func persianTagNamesEvaluator(fld reflect.StructField) string {
	// fieldLabels maps DTO fields json tag
	fieldLabels := map[string]string{
		"order":            "ترتیب",
		"sort":             "مرتب سازی",
		"uuid":             "شناسه",
		"username":         "نام کاربری",
		"password":         "رمز عبور",
		"firstname":        "نام",
		"lastname":         "نام خانوادگی",
		"mobile":           "موبایل",
		"email":            "ایمیل",
		"token":            "توکن",
		"refreshToken":     "رفرش توکن",
		"currentPassword":  "رمز فعلی",
		"newPassword":      "رمز جدید",
		"confirmation":     "تکرار رمز جدید",
		"client":           "مشتری",
		"clientType":       "نوع مشتری",
		"nameFamily":       "نام و نام خانوادگی",
		"description":      "توضیحات",
		"type":             "نوع",
		"active":           "وضعیت فعال",
		"permission":       "مجوز",
		"resource":         "منبع",
		"entity":           "موجودیت",
		"action":           "اقدام",
		"target":           "نشان",
		"role":             "نقش",
		"roleUuid":         "شناسه نقش",
		"permissionsIds":   "شناسه مجوزها",
		"address":          "آدرس",
		"search":           "جستجو",
		"dealType":         "نوع معامله",
		"property":         "ملک",
		"county":           "منطقه",
		"landArea":         "مساحت ملک",
		"minLandArea":      "حداقل مساحت ملک",
		"maxLandArea":      "حداکثر مساحت ملک",
		"buildingArea":     "مساحت سازه",
		"minBuildingArea":  "حداقل مساحت سازه",
		"maxBuildingArea":  "حداکثر مساحت سازه",
		"buildingState":    "وضعیت سازه",
		"buildingAge":      "سن سازه",
		"minBuildingAge":   "حداقل سن سازه",
		"maxBuildingAge":   "حداکثر سن سازه",
		"latitude":         "عرض جغرافیایی",
		"longitude":        "طول جغرافیایی",
		"titleDeed":        "سند مالکیت",
		"titleDeedStatus":  "وضعیت سند مالکیت",
		"areaType":         "محدوده ملکی",
		"constructionArea": "داخل بافت",
		"propertyUse":      "کاربری ملک",
		"features":         "امکانات",
		"price":            "قیمت",
		"minPrice":         "حداقل قیمت",
		"maxPrice":         "حداکثر قیمت",
		"rentPrice":        "اجاره",
		"minRentPrice":     "حداقل اجاره",
		"maxRentPrice":     "حداکثر اجاره",
		"visited":          "بازدید",
		"status":           "وضعیت",
		"share":            "وضعیت اشتراک گذاری",
		"release":          "وضعیت انتشار",
		"releaseDate":      "تاریخ انتشار",
		"expireDate":       "تاریخ پایان",
	}

	jsonName := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if faLabel, exists := fieldLabels[jsonName]; exists {
		return faLabel
	}

	return fld.Name
}
