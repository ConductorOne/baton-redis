package client

import "time"

type User struct {
	AuthMethod             string        `json:"auth_method"`
	BdbsEmailAlerts        []interface{} `json:"bdbs_email_alerts"`
	CertificateSubjectLine string        `json:"certificate_subject_line"`
	ClusterEmailAlerts     bool          `json:"cluster_email_alerts"`
	Email                  string        `json:"email"`
	EmailAlerts            bool          `json:"email_alerts"`
	Name                   string        `json:"name"`
	PasswordIssueDate      time.Time     `json:"password_issue_date"`
	Role                   string        `json:"role"`
	RoleUIDs               []int         `json:"role_uids"`
	Status                 string        `json:"status"`
	UID                    int           `json:"uid"`
}

type Role struct {
	Management string `json:"management"`
	Name       string `json:"name"`
	UID        int    `json:"uid"`
}
