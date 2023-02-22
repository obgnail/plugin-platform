package common

import "fmt"

type IVersion interface {
	Major() int
	Minor() int
	Revision() int
	VersionString() string
	Relationship(IVersion) VersionTime
}

type Version struct {
	major    int
	minor    int
	revision int
}

func NewVersion(major, minor, revision int) IVersion {
	return &Version{
		major:    major,
		minor:    minor,
		revision: revision,
	}
}

func (v *Version) Major() int {
	return v.major
}

func (v *Version) Minor() int {
	return v.minor
}

func (v *Version) Revision() int {
	return v.revision
}

func (v *Version) VersionString() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.revision)
}

func (v *Version) Relationship(t IVersion) VersionTime {
	if v.major > t.Major() || (v.major == t.Major() && v.minor > t.Minor()) ||
		(v.major == t.Major() && v.minor == t.Minor() && v.revision > t.Revision()) {
		return Later
	} else if v.major == t.Major() && v.minor == t.Minor() && v.revision == t.Revision() {
		return Same
	} else {
		return Earlier
	}
}

func (v *Version) GetVersion() IVersion {
	var ret IVersion
	ret = v
	return ret
}

func ParseVersionString(vs string) (*Version, error) {
	var major, minor, revision int
	_, err := fmt.Sscanf(vs, "%d.%d.%d", &major, &minor, &revision)
	if err != nil {
		return &Version{}, err
	}
	version := &Version{
		major:    major,
		minor:    minor,
		revision: revision,
	}
	return version, nil
}
