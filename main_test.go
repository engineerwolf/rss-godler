package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func Test_extractAtom(t *testing.T) {
	type args struct {
		feedInfo gofeed.Item
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Content-dispositionMissingTitleWithoutExt",
			args: args{
				feedInfo: gofeed.Item{
					Title: "Sometitle",
					GUID:  "http://www.google.com",
				},
			},
			want:  "Sometitle.torrent",
			want1: "http://www.google.com",
		},
		{
			name: "Content-dispositionMissingTitleWithExt",
			args: args{
				feedInfo: gofeed.Item{
					Title: "Sometitle.torrent",
					GUID:  "http://www.google.com",
				},
			},
			want:  "Sometitle.torrent",
			want1: "http://www.google.com",
		},
		{
			name: "Content-dispositionMissingTitleMagnetExt",
			args: args{
				feedInfo: gofeed.Item{
					Title: "Sometitle.magnet",
					GUID:  "http://www.google.com",
				},
			},
			want:  "Sometitle.magnet",
			want1: "http://www.google.com",
		},
		{
			name: "Content-dispositionWithExt",
			args: args{
				feedInfo: gofeed.Item{
					Title: "Sometitle",
					GUID:  "https://nyaa.si/view/1196343/torrent",
				},
			},
			want:  "[HorribleSubs] Ore wo Suki nano wa Omae dake ka yo - 08 [720p].mkv.torrent",
			want1: "https://nyaa.si/view/1196343/torrent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractAtom(tt.args.feedInfo)
			if got != tt.want {
				t.Errorf("extractAtom() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractAtom() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
