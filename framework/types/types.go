package types

import "context"

type Aggregate interface {
	ReactTo(Event) error
}

type AggregateItterator interface {
	Len() (int, bool)
	Next() Aggregate
}

// Event interface may be any type which may carry any baggage it likes.
//
// It must serialize and deserialize cleanly for storage reasons.
type Event interface{}

// EventItterator is a simple Itterator interface which should mean that we
// are never in a situation where an aggregate with a large number of
// events causes massive allocations or other resource starvation when
// being rehydrated.
//
// Implementations may choose to batch their reads to the underlaying
// storage into bulk, and iterate over single items at the API level, this
// has been shown to be very performant when talking to Redis for example
// in "pages" of 1000.
//
// Beware that failing to read some itterator implementations to the end
// may hold locks on some underlaying resources.
type EventItterator interface {
	Len() (int, bool)
	Next() Event
}

// Depot is a general storage interface for application related data. The
// name is chosen to draw parallels with the Repository concept sometimes
// used in CQRS/Event Sourcing applications.
//
// The name is chosen because in Retro the Depot also stores metrics about
// the commands themselves, as well as performance information, not simply
// the Events generated by commands.
//
// This is not so much "Event Sourcing" as it is Command/Event sourcing.
type Depot interface {
	Rehydrate(context.Context, Aggregate, string) error
	GetByDirname(context.Context, string) AggregateItterator
}

// Logger is the generic logging interface. It explicitly avoids including
// Fatal and Fatalf because of the relative brutal nature of os.Exit
// without a chance to clean up.
//
// In general tracing should be preferred to logging, however logging can
// always be valuable.
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type CommandDesc interface {
	Name() string
	Path() string
	Args() ApplicationCmdArgs
}

type ApplicationCmdArgs interface{}

type StateEngine interface {
	Apply(context.Context, SessionID, CommandDesc) (string, error)
	StartSession(SessionParams)
}

type SessionID string

type SessionParams interface{}

type AggregateManifest interface {
	Register(string, Aggregate) error
	ForPath(string) (Aggregate, error)
}

type EventManifest interface {
	Register(Event) error
}

// CommandManifest is the interface that allows for various implementations
// of mapping command types to aggregates. They are stored internally in a
// map of types to a slice of strings.
//
// Register takes an aggregate and derives it's type and appends the
// Command type to a slice. There is room for alternative implementations
// which may be faster, or do a smarter search than the range loop to find
// matching commands.
//
// ForAggregate is counterpart to Register, it returns the Commands ready
// to apply, or an error.
type CommandManifest interface {
	Register(Aggregate, Command) error
	ForAggregate(Aggregate) ([]Command, error)
}

// Command is the generic interface to express a user intent towards a
// model in the system.
//
// Commands exist to carry state, the primary calling method is to pass a
// reference to the Apply() function to the calling site, our public
// interface then simply demands that we can pass a simple function, not
// the entire object (closures ensure that the object context is always
// available)
//
// SetState is used to infuse the command with Aggregate state to whom it
// is attached. The aggregate state is embedded into the struct
// implementing Command rathe than given as an argument to express that
// logically one is calling a method *on* an aggregate that has been
// brought upto a certain state.
//
// Apply takes a context, an Aggregate which is expected to represent the
// current session, and a Depot which it may use to look up any other
// Aggregates that it needs to apply business logic.
type Command interface {
	SetState(Aggregate) error
	Apply(context.Context, Aggregate, Depot) ([]Event, error)
}

type CommandFunc func(context.Context, Aggregate, Depot) ([]Event, error)

type Resolver interface {
	Resolve(context.Context, Depot, []byte) (CommandFunc, error)
}

type ResolveFunc func(context.Context, Depot, []byte) (CommandFunc, error)
