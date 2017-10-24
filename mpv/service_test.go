package mpv_test

import (
	"testing"
	"time"
)

// TestService_PlayVideo tests playing a video.
func TestService_PlayVideo(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)
}

// TestService_MultipleVideos_Sequence tests playing two videos in sequence.
func TestService_MultipleVideos_Sequence(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)
}

// TestService_MultipleVideos_Interrupt tests playing two videos replacing the first.
func TestService_MultipleVideos_Interrupt(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(5 * time.Second)

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)
}

// TestService_MultipleVideos_Delay tests playing two videos in sequence with delay.
func TestService_MultipleVideos_Delay(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(12 * time.Second)

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)
}

// TestService_Restart tests restarting the player.
func TestService_Restart(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(4 * time.Second)
	c.Close()

	if err := c.Client.Open(); err != nil {
		panic(err)
	}

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(4 * time.Second)
}

// TestService_ServiceRestart tests restarting the player through the service.
func TestService_ServiceRestart(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(4 * time.Second)

	// stop the player.
	if err := s.Stop(); err != nil {
		panic(err)
	}

	// start the player.
	if err := s.Start(); err != nil {
		panic(err)
	}

	if err := s.PlayFile("../test/test.mp4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(4 * time.Second)
}

// TestService_PlayYoutube tests playing a youtube video.
func TestService_PlayYoutube(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	if err := s.PlayFile("https://www.youtube.com/watch?v=Zz4FLjMEKf4"); err != nil {
		t.Fatal("failed to play video")
	}
	time.Sleep(10 * time.Second)
}
