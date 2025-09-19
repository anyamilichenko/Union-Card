package code

var (
	Success = ResultCode{
		Code:    200,
		Message: "Success",
	}
	BadRequest = ResultCode{
		Code:    400,
		Message: "Bad request",
	}
	Unauthorized = ResultCode{
		Code:    401,
		Message: "You are not authorized",
	}
	Forbidden = ResultCode{
		Code:    403,
		Message: "Permission denied",
	}
	UnprocessableEntity = ResultCode{
		Code:    422,
		Message: "Unprocessable entity",
	}
	InternalServerError = ResultCode{
		Code:    500,
		Message: "Internal server error",
	}
	UserDoesNotExist = ResultCode{
		Code:    1,
		Message: "User does not exist",
	}
	ProfileDoesNotExist = ResultCode{
		Code:    2,
		Message: "Profile does not exist",
	}
	InvalidPassword = ResultCode{
		Code:    3,
		Message: "Invalid password",
	}
	UserAlreadyExists = ResultCode{
		Code:    4,
		Message: "User already exists",
	}
	EmailIsBusy = ResultCode{
		Code:    5,
		Message: "Email is busy",
	}
	EmailIsFree = ResultCode{
		Code:    6,
		Message: "Email is free",
	}
	UserPasswordIsNotSet = ResultCode{
		Code:    7,
		Message: "The password is not set, login through third-party service",
	}
	SocialDoesNotExist = ResultCode{
		Code:    8,
		Message: "Social does not exist",
	}
)
