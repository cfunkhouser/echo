package echo

import (
	"os"
	"slices"
	"strings"
)

// EnvVar is an environment variable value stored in memory as a name / value
// pair.
type EnvVar struct {
	Name, Value string
}

// EnvironmentDump represents the bits of the running environment to be dumped
// for return in a response payload.
type EnvironmentDump struct {
	Variables []*EnvVar
}

// EnvironmentDumpPolicy determines what should be dumped from the running
// process' environment.
type EnvironmentDumpPolicy func(*EnvVar) bool

// denyAll prevents the environment from being dumped.
func denyAll(_ *EnvVar) bool {
	return false
}

// VeryDangerousAllowAll is an EnvironmentDumpPolicy which dumps the entire
// environment. Aim this away from your foot before using.
func VeryDangerousAllowAll(_ *EnvVar) bool {
	return true
}

const allowedVarsVarName = "ECHO_VARS"

// PolicyFromEnv is an EnvironmentDumpPolicy which reads the value of the
// ECHO_VARS environment variable, and allows dumping any env var names found.
// If the value is set to "..." then the entire environment is dumped (which is
// often a bad idea).
func PolicyFromEnv() func(*EnvVar) bool {
	allowed := make(map[string]bool)
	vv, ok := os.LookupEnv(allowedVarsVarName)
	if !ok {
		return denyAll
	}

	vv = strings.TrimSpace(vv)
	if vv == "..." {
		return VeryDangerousAllowAll
	}

	for _, vn := range strings.Split(vv, " ") {
		vn = strings.TrimSpace(vn)
		if vn == "" {
			continue
		}
		allowed[vn] = true
	}
	return func(v *EnvVar) bool {
		return allowed[v.Name]
	}
}

var defaultPolicy = PolicyFromEnv()

// DumpEnv for rendering in an HTML payload.
func DumpEnv(policy EnvironmentDumpPolicy) *EnvironmentDump {
	if policy == nil {
		policy = defaultPolicy
	}

	var vars []*EnvVar
	for _, rv := range os.Environ() {
		vs := strings.SplitN(rv, "=", 2)
		vn, vv := strings.TrimSpace(vs[0]), ""
		if len(vs) > 1 {
			vv = strings.TrimSpace(vs[1])
		}
		v := &EnvVar{vn, vv}
		if policy(v) {
			vars = append(vars, &EnvVar{vn, vv})
		}
	}

	// Sort the vars for determinism and human friendliness.
	slices.SortFunc(vars, func(l, r *EnvVar) int {
		return strings.Compare(l.Name, r.Name)
	})

	return &EnvironmentDump{
		Variables: vars,
	}
}
