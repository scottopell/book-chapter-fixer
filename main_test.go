package main

import "testing"

func Test_getChapter(t *testing.T) {
	type args struct {
		fileName string
		tagTitle string
	}
	tests := []struct {
		name        string
		args        args
		newFileName string
		newTitle    string
		origChapNum int
		newChapNum  int
		wantErr     bool
	}{
		{"High chapter",
			args{"Book Name (Chapter 67).mp3", "Book Name - 67"},
			"Book Name (Chapter 63).mp3",
			"Book Name - Chapter 63",
			67,
			63,
			false,
		},
		{"low chapter",
			args{"Book Name (Chapter 01).mp3", "Book Name - 01"},
			"Book Name (Chapter -3).mp3",
			"Book Name - Chapter -3",
			1,
			-3,
			false,
		},
		{"mismatch",
			args{"Book Name (Chapter 01).mp3", "Book Name - 04"},
			"",
			"",
			0,
			0,
			true,
		},
		{"wrong filename format",
			args{"Book Name Chapter 01.mp3", "Book Name - 04"},
			"",
			"",
			0,
			0,
			true,
		},
		{"wrong title format",
			args{"Book Name (Chapter 01).mp3", "Book Name 04"},
			"",
			"",
			0,
			0,
			true,
		},
		{"Different Prefixes",
			args{"Booooook Name (Chapter 57).mp3", "Book Name - 57"},
			"Booooook Name (Chapter 53).mp3",
			"Book Name - Chapter 53",
			57,
			53,
			false,
		},
		{"Different Prefixes 2",
			args{"Booooook Name (Chapter 40).mp3", "Bok Bok Name - 40"},
			"Booooook Name (Chapter 36).mp3",
			"Bok Bok Name - Chapter 36",
			40,
			36,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filenamePrefix, titlePrefix, origChap, fixedChap, err := getChapter(tt.args.fileName, tt.args.tagTitle)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if filenamePrefix != tt.newFileName {
				t.Errorf("getChapter() got = %v, want %v", filenamePrefix, tt.newFileName)
			}
			if titlePrefix != tt.newTitle {
				t.Errorf("getChapter() got1 = %v, want %v", titlePrefix, tt.newTitle)
			}
			if origChap != tt.origChapNum {
				t.Errorf("getChapter() got2 = %v, want %v", origChap, tt.origChapNum)
			}
			if fixedChap != tt.newChapNum {
				t.Errorf("getChapter() got3 = %v, want %v", fixedChap, tt.newChapNum)
			}
		})
	}
}
