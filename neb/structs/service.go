package structs

/*

#cgo nagios3 CFLAGS: -DNAGIOS3 -I. -I${SRCDIR}/../../libs
#cgo nagios3 LDFLAGS: -Wl,-unresolved-symbols=ignore-all

#cgo nagios4 CFLAGS: -DNAGIOS4 -I. -I${SRCDIR}/../../libs
#cgo nagios4 LDFLAGS: -Wl,-unresolved-symbols=ignore-all

#cgo naemon CFLAGS: -DNAEMON -I.
#cgo naemon pkg-config: naemon

#include "../dependencies.h"

//This block is needed due to different naming schemes in Nagios3 and Nagios4/Naemon
#if defined(NAGIOS3)
char* ServiceGetCommand(void* data) {
	return ((service *)data)->service_check_command;
}
#elif defined(NAGIOS4) || defined(NAEMON)
char* ServiceGetCommand(void* data) {
	return ((service *)data)->check_command;
}
#endif

*/
import "C"
import (
	"strings"
	"unsafe"
)

//Servicelist is a list of services
type Servicelist []Service

type Service struct {
	HostName    string
	Description string
	DisplayName string
	//CheckCommand contains args
	CheckCommand string
	//Command is the pure pluginname
	Command       string
	ChecksEnabled int
	CheckType     int
	IsFlapping    int
}

//CastService tries to cast the pointer to an go struct
func CastService(data unsafe.Pointer) Service {
	st := *((*C.service)(data))
	command := C.GoString(C.ServiceGetCommand(data))
	return Service{
		HostName:      C.GoString(st.host_name),
		Description:   C.GoString(st.description),
		DisplayName:   C.GoString(st.display_name),
		CheckCommand:  command,
		Command:       splitCommand(command),
		ChecksEnabled: int(st.checks_enabled),
		CheckType:     int(st.check_type),
		IsFlapping:    int(st.is_flapping),
	}
}

func splitCommand(checkCommand string) string {
	if strings.Contains(checkCommand, "!") {
		return strings.Split(checkCommand, "!")[0]
	} else {
		return checkCommand
	}
}
