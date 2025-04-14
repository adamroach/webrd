#import <AVFoundation/AVFoundation.h>
#import <CoreMedia/CoreMedia.h>
#import <CoreVideo/CoreVideo.h>
#import <Foundation/Foundation.h>

@interface Capturer : NSObject <AVCaptureVideoDataOutputSampleBufferDelegate> {
  @private
    AVCaptureSession *mSession;
    AVCaptureVideoDataOutput *mVideoDataOutput;
    dispatch_queue_t mVideoDataOutputQueue;
    void *mCallbackOpaque;
}

- (void)start:(void *)opaque;
- (void)stop;

@end