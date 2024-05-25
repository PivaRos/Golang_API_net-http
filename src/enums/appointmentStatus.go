package enums

type AppointmentStatus string

const (
	Pending  AppointmentStatus = "Pending"
	Approved AppointmentStatus = "Approved"
	Removed  AppointmentStatus = "Removed"
)
