package ltm

type Status int8

func (s Status) String() string {
	if s > 21 {
		return "STATUS_UNKNOWN"
	}

	return []string{
		"STATUS_MANUAL",
		"STATUS_RATE",
		"STATUS_ANGLE",
		"STATUS_HORIZON",
		"STATUS_ACRO",
		"STATUS_STABILISED1",
		"STATUS_STABILISED2",
		"STATUS_STABILISED3",
		"STATUS_ALTITUDE_HOLD",
		"STATUS_GPS_HOLD",
		"STATUS_WAYPOINTS",
		"STATUS_HEAD_FREE",
		"STATUS_CIRCLE",
		"STATUS_RTH",
		"STATUS_FOLLOW_ME",
		"STATUS_LAND",
		"STATUS_FLY_BY_WIREA",
		"STATUS_FLY_BY_WIREB",
		"STATUS_CRUISE",
		"STATUS_UNKNOWN",
		"STATUS_LAUNCH",
		"STATUS_AUTOTUNE",
	}[s]
}

type GPSMode int8

func (m GPSMode) String() string {
	if m > 3 {
		return "UNKNOWN"
	}

	return []string{
		"NONE", "POSHOLD", "RTH", "MISSION",
	}[m]

}

type NavMode int8

func (m NavMode) String() string {
	if m > 15 {
		return "UNKNOWN"
	}

	return []string{
		"NONE", "RTH_START", "RTH_ENROUTE", "POSHOLD_INF", "POSHOLD_TIMED",
		"WP_ENROUTE", "PROCESS_NEXT", "JUMP", "START_LAND", "LANDING_INPROGRESS",
		"LANDED", "SETTLING_BEFORE_LANDING", "START_DESCENT", "HOVER_ABOVE_HOME",
		"EMERGENCY_LANDING", "CRITICAL_GPS_FAILURE",
	}[m]
}

type NavAction int8

func (a NavAction) String() string {
	if a > 8 {
		return "UNKNOWN"
	}

	return []string{
		"UNASSIGNED", "WAYPOINT", "POSHOLD_UNLIM", "POSHOLD_TIME", "RTH",
		"SET_POI", "JUMP", "SET_HEAD", "LAND",
	}[a]
}

type NavError int8

func (e NavError) String() string {
	if e > 11 {
		return "UNKNOWN"
	}

	return []string{
		"Navigation system is working",
		"Next waypoint distance is more than the safety limit, aborting mission",
		"GPS reception is compromised - pausing mission",
		"Error while reading next waypoint from memory, aborting mission",
		"Mission Finished",
		"Waiting for timed position hold",
		"Invalid Jump target detected, aborting mission",
		"Invalid Mission Step Action code detected, aborting mission",
		"Waiting to reach return to home altitude",
		"GPS fix lost, mission aborted",
		"Disarmed, navigation engine disabled",
		"Landing is in progress, check attitude",
	}[e]
}
