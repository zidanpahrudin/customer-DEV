package dto

import "time"

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// CustomerRequest represents customer creation/update request
type CustomerRequest struct {
	Name string `json:"name" binding:"required" example:"PT Teknologi Maju"`
	/* Email       string  `json:"email" example:"info@teknologimaju.com"` */
	/* Phone       string  `json:"phone" example:"021-12345678"`
	Website     string  `json:"website" example:"https://teknologimaju.com"` */
	/* Description string  `json:"description" example:"Perusahaan teknologi informasi"` */
	Status      string  `json:"status" example:"Active"`
	Category    string  `json:"category" example:"Technology"`
	Rating      float64 `json:"rating" example:"4.5"`
	AverageCost float64 `json:"average_cost" example:"50000000"`
}

// CustomersResponse represents customers list response
type CustomersResponse struct {
	Customers []Customer `json:"customers"`
	Stats     Stats      `json:"stats"`
}

// Stats represents customer statistics
type Stats struct {
	TotalCustomers   int64   `json:"total_customers" example:"100"`
	NewCustomers     int64   `json:"new_customers" example:"10"`
	AvgCost          float64 `json:"avg_cost" example:"45000000"`
	BlockedCustomers int64   `json:"blocked_customers" example:"5"`
}

// User represents user data in responses
type User struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"user123"`
	Email    string `json:"email" example:"user@example.com"`
	RoleID   uint   `json:"role_id" example:"1"`
}

// Customer represents customer data
type Customer struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"PT Teknologi Maju"`
	BrandName string `json:"brand_name" example:"TechMaju"`
	Code      string `json:"code" example:"TM001"`
	/* Email       string  `json:"email" example:"info@teknologimaju.com"`
	Phone       string  `json:"phone" example:"021-12345678"`
	Website     string  `json:"website" example:"https://teknologimaju.com"` */
	/* Description string  `json:"description" example:"Perusahaan teknologi informasi"` */
	Logo        string  `json:"logo" example:"uploads/logos/logo_1.png"`
	LogoSmall   string  `json:"logo_small" example:"uploads/logos_small/logo_small_1.png"` // Field baru untuk logo kecil
	Status      string  `json:"status" example:"Active"`
	Category    string  `json:"category" example:"Technology"`
	Rating      float64 `json:"rating" example:"4.5"`
	AverageCost float64 `json:"average_cost" example:"50000000"`
}

// CustomerResponse represents customer with simplified relations
type CustomerResponse struct {
	ID               uint   `json:"id" example:"1"`
	Name             string `json:"name" example:"PT Teknologi Maju"`
	BrandName        string `json:"brand_name" example:"TechMaju"`
	Code             string `json:"code" example:"TM001"`
	AccountManagerId string `json:"account_manager_id" example:"1"`
	/* Email            string            `json:"email" example:"info@teknologimaju.com"`
	Phone            string            `json:"phone" example:"021-12345678"`
	Website          string            `json:"website" example:"https://teknologimaju.com"`
	Description      string            `json:"description" example:"Perusahaan teknologi informasi"` */
	Logo        string            `json:"logo" example:"uploads/logos/logo_1.png"`
	LogoSmall   string            `json:"logo_small" example:"uploads/logos_small/logo_small_1.png"`
	Status      string            `json:"status" example:"Active"`
	Category    string            `json:"category" example:"Technology"`
	Rating      float64           `json:"rating" example:"4.5"`
	AverageCost float64           `json:"average_cost" example:"50000000"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Addresses   []AddressResponse `json:"addresses,omitempty"`
	Contacts    []ContactResponse `json:"contacts,omitempty"`
	Others      []OtherResponse   `json:"others,omitempty"`
}

// CreateCustomerRequest represents comprehensive customer creation request
type CreateCustomerRequest struct {
	Name             string                   `json:"name" binding:"required"`
	BrandName        string                   `json:"brandName"`
	Code             string                   `json:"code"`
	AccountManagerId string                   `json:"accountManagerId"`
	Logo             *string                  `json:"logo"`
	LogoSmall        *string                  `json:"logoSmall"`
	StatusName       string                   `json:"status_name"`
	Addresses        []CreateAddressRequest   `json:"addresses,omitempty"`
	Socials          []CreateSocialRequest    `json:"socials,omitempty"`
	Contacts         []CreateContactRequest   `json:"contacts,omitempty"`
	Structures       []CreateStructureRequest `json:"structures,omitempty"`
	Groups           CreateGroupsRequest      `json:"groups,omitempty"` // Ubah dari []CreateGroupsRequest ke CreateGroupsRequest
	Others           []CreateOtherRequest     `json:"others,omitempty"`
}

// CreateAddressRequest represents address creation in customer request
type CreateAddressRequest struct {
	// CustomerID uint   `json:"customer_id" binding:"required"` // Hapus field ini
	Name    string `json:"name" binding:"required" example:"Head Office"`
	Address string `json:"address" binding:"required" example:"Jl. Sudirman No. 123, Jakarta Selatan"`
	IsMain  bool   `json:"isMain" example:"true"`
	Active  bool   `json:"active" example:"true"`
}

// CreateSocialRequest represents social media creation in customer request
type CreateSocialRequest struct {
	// Name     string `json:"name" binding:"required" example:"Instagram"` // Hapus field ini karena duplikat dengan Platform
	Platform string `json:"platform" binding:"required" example:"Instagram"`
	Handle   string `json:"handle" binding:"required" example:"@digiinno_id"`
	Active   bool   `json:"active" example:"true"`
}

// CreateContactRequest represents contact creation in customer request
type CreateContactRequest struct {
	// CustomerID  uint   `json:"customer_id" binding:"required"` // Hapus field ini
	Name        string `json:"name" binding:"required" example:"Budi Santoso"`
	Birthdate   string `json:"birthdate" example:"1985-03-15"`
	JobPosition string `json:"jobPosition" example:"CEO"`
	Email       string `json:"email" example:"budi@digiinno.com"`
	Phone       string `json:"phone" example:"021-5551234"`
	Mobile      string `json:"mobile" example:"0812-3456-7890"`
	IsMain      bool   `json:"isMain" example:"true"`
	Active      bool   `json:"active" example:"true"`
}

// CreateStructureRequest represents structure creation in customer request
type CreateStructureRequest struct {
	// CustomerID uint    `json:"customer_id" binding:"required"` // Hapus field ini
	TempKey   string  `json:"tempKey" example:"1"`
	ParentKey *string `json:"parentKey" example:"null"`
	Name      string  `json:"name" binding:"required" example:"Board of Directors"`
	Level     int     `json:"level" binding:"required" example:"1"`
	Address   string  `json:"address" example:"Jakarta"`
	Active    bool    `json:"active" example:"true"`
}

// CreateGroupsRequest represents groups assignment in customer request
type CreateGroupsRequest struct {
	IndustryID        string `json:"industryId" example:"1"` // Perbaiki nama field
	IndustryActive    bool   `json:"industryActive" example:"true"`
	ParentGroupID     string `json:"parentGroupId" example:"2"` // Perbaiki nama field
	ParentGroupActive bool   `json:"parentGroupActive" example:"true"`
}

// CreateOtherRequest represents other attributes in customer request
type CreateOtherRequest struct {
	// CustomerID uint    `json:"customer_id" binding:"required"` // Hapus field ini
	Key    string  `json:"key" binding:"required" example:"company_size"`
	Value  *string `json:"value" example:"50-100 employees"`
	Active bool    `json:"active" example:"true"`
}

// CreateActivityRequest represents activity creation request
type CreateActivityRequest struct {
	CustomerID   uint   `json:"customer_id" binding:"required" example:"1"`
	Title        string `json:"title" binding:"required" example:"Client Meeting"`
	Type         string `json:"type" binding:"required" example:"Meeting"`
	Agenda       string `json:"agenda" example:"Discuss project requirements"`
	StartTime    string `json:"start_time" binding:"required" example:"2024-01-15T10:00:00Z"`
	EndTime      string `json:"end_time" binding:"required" example:"2024-01-15T12:00:00Z"`
	LocationName string `json:"location_name" example:"Conference Room A"`
	Status       string `json:"status" example:"Scheduled"`
	// Hapus field Lat dan Lng yang masih ada
}

// UpdateActivityRequest represents activity update request
type UpdateActivityRequest struct {
	Title        *string `json:"title" example:"Updated Meeting"`
	Type         *string `json:"type" example:"Meeting"`
	Agenda       *string `json:"agenda" example:"Updated agenda"`
	StartTime    *string `json:"start_time" example:"2024-01-15T10:00:00Z"`
	EndTime      *string `json:"end_time" example:"2024-01-15T12:00:00Z"`
	LocationName *string `json:"location_name" example:"Conference Room B"`
	Status       *string `json:"status" example:"Completed"`
}

// ActivityResponse represents activity response
type ActivityResponse struct {
	ID           uint   `json:"id" example:"1"`
	CustomerID   uint   `json:"customer_id" example:"1"`
	Title        string `json:"title" example:"Client Meeting"`
	Type         string `json:"type" example:"Meeting"`
	Agenda       string `json:"agenda" example:"Discuss project requirements"`
	StartTime    string `json:"start_time" example:"2024-01-15T10:00:00Z"`
	EndTime      string `json:"end_time" example:"2024-01-15T12:00:00Z"`
	LocationName string `json:"location_name" example:"Conference Room A"`
	Status       string `json:"status" example:"Scheduled"`
	CreatedBy    uint   `json:"created_by" example:"1"`
	CreatedAt    string `json:"created_at" example:"2024-01-15T08:00:00Z"`
	UpdatedAt    string `json:"updated_at" example:"2024-01-15T08:00:00Z"`
}

// ActivitiesResponse represents activities list response
type ActivitiesResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total" example:"10"`
}

// ActivityAttendeeRequest represents activity attendee request
type ActivityAttendeeRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required" example:"[1,2,3]"`
}

// ActivityCheckinRequest represents activity check-in request
type ActivityCheckinRequest struct {
	// Bisa ditambahkan field jika diperlukan, misalnya:
	// Notes string `json:"notes" example:"Arrived on time"`
	// Location string `json:"location" example:"Conference Room A"`
}

// Invoice DTOs
/* type CreateInvoiceRequest struct {
	CustomerID    uint      `json:"customer_id" binding:"required"`
	ProjectID     string    `json:"project_id"`
	InvoiceNumber string    `json:"invoice_number" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	IssuedDate    time.Time `json:"issued_date" binding:"required"`
	DueDate       time.Time `json:"due_date" binding:"required"`
	PaidAmount    float64   `json:"paid_amount" binding:"gte=0"`
} */

/* type UpdateInvoiceRequest struct {
	ProjectID     *string    `json:"project_id"`
	InvoiceNumber *string    `json:"invoice_number"`
	Amount        *float64   `json:"amount" binding:"omitempty,gt=0"`
	IssuedDate    *time.Time `json:"issued_date"`
	DueDate       *time.Time `json:"due_date"`
	PaidAmount    *float64   `json:"paid_amount" binding:"omitempty,gte=0"`
} */

/* type InvoiceResponse struct {
	ID            uint              `json:"id"`
	CustomerID    uint              `json:"customer_id"`
	ProjectID     string            `json:"project_id"`
	InvoiceNumber string            `json:"invoice_number"`
	Amount        float64           `json:"amount"`
	IssuedDate    time.Time         `json:"issued_date"`
	DueDate       time.Time         `json:"due_date"`
	PaidAmount    float64           `json:"paid_amount"`
	Balance       float64           `json:"balance"` // Amount - PaidAmount
	Status        string            `json:"status"`  // Paid, Partial, Unpaid, Overdue
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Customer      *CustomerResponse `json:"customer,omitempty"`
} */

// Payment DTOs
/* type CreatePaymentRequest struct {
	InvoiceID uint      `json:"invoice_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required,gt=0"`
	PaidAt    time.Time `json:"paid_at" binding:"required"`
}

type UpdatePaymentRequest struct {
	Amount *float64   `json:"amount" binding:"omitempty,gt=0"`
	PaidAt *time.Time `json:"paid_at"`
} */

/* type PaymentResponse struct {
	ID        uint             `json:"id"`
	InvoiceID uint             `json:"invoice_id"`
	Amount    float64          `json:"amount"`
	PaidAt    time.Time        `json:"paid_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Invoice   *InvoiceResponse `json:"invoice,omitempty"`
} */

// Status DTOs
type CreateStatusRequest struct {
	StatusName string `json:"status_name" binding:"required"`
}

type UpdateStatusRequest struct {
	StatusName string `json:"status_name" binding:"required"`
}

type StatusResponse struct {
	ID         uint      `json:"id"`
	StatusName string    `json:"status_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AddressResponse represents simplified address response
type AddressResponse struct {
	Name    string `json:"name" example:"Kantor Utama"`
	Address string `json:"address" example:"Jl. Merdeka No. 45, Surabaya"`
	IsMain  bool   `json:"isMain" example:"true"`
	Active  bool   `json:"active" example:"true"`
}

// ContactResponse represents simplified contact response
type ContactResponse struct {
	Name        string `json:"name" example:"Bambang Sutrisno"`
	JobPosition string `json:"jobPosition" example:"Owner"`
	Email       string `json:"email" example:"bambang@berkahjaya.com"`
	Phone       string `json:"phone" example:"031-1122334"`
	Mobile      string `json:"mobile" example:"0856-7788-9900"`
	IsMain      bool   `json:"isMain" example:"true"`
	Active      bool   `json:"active" example:"true"`
}

// OtherResponse represents simplified other response
type OtherResponse struct {
	Key    string `json:"key" example:"company_size"`
	Value  string `json:"value" example:"10-25 employees"`
	Active bool   `json:"active" example:"true"`
}
