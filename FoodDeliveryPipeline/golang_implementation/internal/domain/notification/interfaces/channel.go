package interfaces

type Channel interface {
	Notify(ctx MessageContext)
}
