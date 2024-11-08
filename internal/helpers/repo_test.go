package helpers

import (
	"testing"
)

func TestExtractUserAndRepo(t *testing.T) {
	tests := []struct {
		name      string
		remoteUrl string
		wantUser  string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "Valid SSH URL",
			remoteUrl: "git@github.com:user/repo.git",
			wantUser:  "user",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Valid HTTPS URL",
			remoteUrl: "https://github.com/user/repo.git",
			wantUser:  "user",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Invalid URL format",
			remoteUrl: "ftp://github.com/user/repo.git",
			wantUser:  "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "Missing .git suffix",
			remoteUrl: "git@github.com:user/repo",
			wantUser:  "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "Valid SSH URL with custom host",
			remoteUrl: "git@custom-host.com:developer/project.git",
			wantUser:  "developer",
			wantRepo:  "project",
			wantErr:   false,
		},
		{
			name:      "Valid HTTPS URL with custom host",
			remoteUrl: "https://custom-host.com/developer/project.git",
			wantUser:  "developer",
			wantRepo:  "project",
			wantErr:   false,
		},
		{
			name:      "Empty URL",
			remoteUrl: "",
			wantUser:  "",
			wantRepo:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotRepo, err := ExtractUserAndRepo(tt.remoteUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUserAndRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUser != tt.wantUser {
				t.Errorf("ExtractUserAndRepo() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("ExtractUserAndRepo() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}
