#import <AVFoundation/AVFoundation.h>
#import <CoreMedia/CoreMedia.h>
#import <CoreVideo/CoreVideo.h>
#import <Foundation/Foundation.h>

@interface VideoCapturer
    : NSObject <AVCaptureVideoDataOutputSampleBufferDelegate> {
  @private
    AVCaptureSession *mSession;
    AVCaptureVideoDataOutput *mVideoDataOutput;
    dispatch_queue_t mVideoDataOutputQueue;
    void *mCallbackOpaque;
}

- (void)start:(void *)opaque fps:(int)fps;
- (void)stop;

@end