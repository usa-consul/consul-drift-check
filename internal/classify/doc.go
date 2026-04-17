// Package classify assigns severity levels (info, warning, critical) to
// drift results based on key-prefix rules. Rules are evaluated in
// declaration order and the first match wins, allowing operators to
// promote sensitive namespaces (e.g. "prod/secrets/") to critical while
// leaving routine configuration keys at the default info level.
package classify
