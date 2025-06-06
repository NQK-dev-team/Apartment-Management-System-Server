package constants

type supportTicketStatusStruct struct {
	PENDING  int
	APPROVED int
	REJECTED int
}

type contractStatusStruct struct {
	ACTIVE                int
	EXPIRED               int
	CANCELLED             int
	WAITING_FOR_SIGNATURE int
	NOT_IN_EFFECT         int
}

type roomStatusStruct struct {
	RENTED      int
	SOLD        int
	AVAILABLE   int
	MAINTENANCE int
	UNAVAILABLE int
}

type userGenderStruct struct {
	MALE   int
	FEMALE int
	OTHER  int
}

type importTypeStruct struct {
	ADD_BUILDINGS int
	ADD_ROOMS     int
	ADD_EMPLOYEES int
	ADD_CUSTOMERS int
	ADD_BILLS     int
}

var Common = struct {
	SupportTicketStatus supportTicketStatusStruct
	ContractStatus      contractStatusStruct
	RoomStatus          roomStatusStruct
	UserGender          userGenderStruct
	ImportType          importTypeStruct
	EmailTokenLimit     int
	NewPasswordLength int
}{
	SupportTicketStatus: supportTicketStatusStruct{
		PENDING:  1,
		APPROVED: 2,
		REJECTED: 3,
	},
	ContractStatus: contractStatusStruct{
		ACTIVE:                1,
		EXPIRED:               2,
		CANCELLED:             3,
		WAITING_FOR_SIGNATURE: 4,
		NOT_IN_EFFECT:         5,
	},
	RoomStatus: roomStatusStruct{
		RENTED:      1,
		SOLD:        2,
		AVAILABLE:   3,
		MAINTENANCE: 4,
		UNAVAILABLE: 5,
	},
	UserGender: userGenderStruct{
		MALE:   1,
		FEMALE: 2,
		OTHER:  3,
	},
	ImportType: importTypeStruct{
		ADD_BUILDINGS: 1,
		ADD_ROOMS:     2,
		ADD_EMPLOYEES: 3,
		ADD_CUSTOMERS: 4,
		ADD_BILLS:     5,
	},
	EmailTokenLimit: 5,
	NewPasswordLength: 8,
}
