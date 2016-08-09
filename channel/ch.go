// An intersting pattern for testing, but you can use the check anywhere
import "testing"

func TestCheckingChannel(t *testing.T) {
  stop := make(chan bool)

  // Testing some fucntion that SHOULD close the channel
  func (stop chan bool) {
    close(chan)
  }(stop)

  // Make sure that the function does close the channel
  _, ok := (<-stop)
  
  // If we can recieve on the channel then it is NOT closed
  if ok {
    t.Error("Channel is not closed")
  }
}
