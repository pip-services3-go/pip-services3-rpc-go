package services


/*
 IRegisterable is interface to perform on-demand registrations.
*/
type IRegisterableT interface {
	// Perform required registration steps.
	Register()
}
