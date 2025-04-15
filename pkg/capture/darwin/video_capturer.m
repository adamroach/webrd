#import "video_capturer.h"
#include <CoreMedia/CMSampleBuffer.h>
#include <CoreVideo/CVImageBuffer.h>
#include <CoreVideo/CVPixelBuffer.h>

extern void process_yuv_frame(void *opaque, void *y, void *cb, void *cr,
                              int yStride, int cStride, int width, int height);

@implementation VideoCapturer
- (void)start:(void *)opaque {
    if (mSession) {
        [self stop];
    }
    mCallbackOpaque = opaque;

    // Create a capture session
    mSession = [[AVCaptureSession alloc] init];

    // Set the session preset as you wish
    mSession.sessionPreset = AVCaptureSessionPresetHigh;

    // If you're on a multi-display system and you want to capture a secondary
    // display, you can call CGGetActiveDisplayList() to get the list of all
    // active displays.
    CGDirectDisplayID displayId = kCGDirectMainDisplay;

    // Create a ScreenInput with the display and add it to the session
    AVCaptureScreenInput *input = [[[AVCaptureScreenInput alloc]
        initWithDisplayID:displayId] autorelease];
    if (!input) {
        [mSession release];
        mSession = nil;
        return;
    }
    input.minFrameDuration = CMTimeMake(1, 30); // 30 fps
    if ([mSession canAddInput:input])
        [mSession addInput:input];

    mVideoDataOutput = [[AVCaptureVideoDataOutput alloc] init];

    // Set up the video capture delegate
    NSDictionary *outputSettings = [NSDictionary
        dictionaryWithObject:
            [NSNumber numberWithInt:kCVPixelFormatType_420YpCbCr8Planar]
                      forKey:(id)kCVPixelBufferPixelFormatTypeKey];
    [mVideoDataOutput setVideoSettings:outputSettings];
    [mVideoDataOutput setAlwaysDiscardsLateVideoFrames:YES];
    mVideoDataOutputQueue =
        dispatch_queue_create("VideoDataOutputQueue", DISPATCH_QUEUE_SERIAL);
    [mVideoDataOutput setSampleBufferDelegate:self queue:mVideoDataOutputQueue];
    [[mVideoDataOutput connectionWithMediaType:AVMediaTypeVideo]
        setEnabled:YES];

    // Add the video capture delegate as an output
    if ([mSession canAddOutput:mVideoDataOutput]) {
        [mSession addOutput:mVideoDataOutput];
    } else {
        NSLog(@"can't add output");
    }

    // Start running the session
    [mSession startRunning];
}

- (void)stop {
    if (mSession) {
        [mSession stopRunning];
        [mSession release];
        mSession = nil;
    }
    if (mVideoDataOutput) {
        [mVideoDataOutput release];
        mVideoDataOutput = nil;
    }
    if (mVideoDataOutputQueue) {
        [mVideoDataOutputQueue release];
        mVideoDataOutputQueue = nil;
    }
    mCallbackOpaque = nil;
}

// AVCaptureVideoDataOutputSampleBufferDelegate methods

- (void)captureOutput:(AVCaptureOutput *)output
    didOutputSampleBuffer:(CMSampleBufferRef)buffer
           fromConnection:(AVCaptureConnection *)connection {

    if (!mCallbackOpaque) {
        return;
    }

    // NOTE: Clients that need to reference the CMSampleBuffer object
    // outside of the scope of this method must CFRetain it and then
    // CFRelease it when they are finished with it.
    CVImageBufferRef img = CMSampleBufferGetImageBuffer(buffer);

    /*
    bool planar = CVPixelBufferIsPlanar(img);
    int planeCount = CVPixelBufferGetPlaneCount(img);
    void *ptr = CVPixelBufferGetBaseAddress(img);
    */

    // The base address must be locked to access any of the image data
    CVPixelBufferLockBaseAddress(img, 0);

    // This is the information Go needs to greate an image.YCbCr
    void *Y = CVPixelBufferGetBaseAddressOfPlane(img, 0);
    void *Cb = CVPixelBufferGetBaseAddressOfPlane(img, 1);
    void *Cr = CVPixelBufferGetBaseAddressOfPlane(img, 2);
    int YStride = CVPixelBufferGetBytesPerRowOfPlane(img, 0);
    int CStride = CVPixelBufferGetBytesPerRowOfPlane(img, 1);
    int Width = CVPixelBufferGetWidth(img);
    int Height = CVPixelBufferGetHeight(img);

    process_yuv_frame(mCallbackOpaque, Y, Cb, Cr, YStride, CStride, Width,
                      Height);

    CVPixelBufferUnlockBaseAddress(img, 0);
}

- (void)captureOutput:(AVCaptureOutput *)output
    didDropOutputSampleBuffer:(CMSampleBufferRef)buffer
               fromConnection:(AVCaptureConnection *)connection {
    NSLog(@"dropped capture");
}
@end