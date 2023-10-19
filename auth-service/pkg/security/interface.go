package security

// TODO: consider using hmac because it simple and satisfied your use case

// ISecurity is interface which provide utility related to web security one of the purpose to make it as an interface is for ease unit testing
type ISecurity interface {
	IPassword
	IGenerator
	IValidator
	ICryptor
}

type ICryptor interface {
	// Decrypt decrypt given encrypted/cipher text with RSA - OAEP algorithm
	Decrypt(chipher string) (string, error)
	EncryptWithURLEncode(plaintext string) (string, error)
}
type IPassword interface {
	// HashPassword wrap for bcrypt.GenerateFromPassword and will generate hash for given password
	HashPassword(password string) (string, error)
	// CompareHashPassword wrap and compare plain text password with given hashed password using bcrypt.CompareHashAndPassword
	CompareHashPassword(password, hashedPassword string) error
}

type IGenerator interface {
	// GenerateSID generate random string for session id purpose
	GenerateSID() (string, error)
	// GenerateToken generate JWT token. it return token as string and error error will be returned as nil if no error occur otherwise it's not nill
	GenerateToken(tokenType string) (string, error)
}

type IValidator interface {
	// ValidateToken validate a JWT token. It returns any/interface{} and nil if no errors occur. In case of an error, it returns an error and nil for the any/interface{} type.
	//
	// The tokenStr argument is expected to be encoded using the standard Base64 encoder.
	ValidateToken(tokenStr string) (any, error)
	ValidateRequest(val any) error
}
