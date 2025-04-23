package auth

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework OpenDirectory -framework Foundation

#import <OpenDirectory/OpenDirectory.h>
#import <Foundation/Foundation.h>

BOOL checkPassword(const char *username, const char *password) {
    NSError *err = nil;
	ODSession *mySession = [ODSession defaultSession];
	ODNode *myNode = [ODNode nodeWithSession:mySession type:kODNodeTypeAuthentication error:&err];
  	if (err) {
  		NSLog(@"Unable to get node: %@", err);
 		return NO;
    }

	NSString *user = [NSString stringWithCString:username encoding:NSUTF8StringEncoding];
	ODRecord *myRecord = [myNode recordWithRecordType:kODRecordTypeUsers
                                                 name:user
                                           attributes:nil
                                                error:&err];
	if (err) {
		NSLog(@"Unable to get %@'s record: %@", user, err);
		return NO;
	}

	NSString *pass = [NSString stringWithCString:password encoding:NSUTF8StringEncoding];
	BOOL retval =  [myRecord verifyPassword:pass error:nil];
	NSLog(@"Password check for user %@: %@", user, (retval ? @"SUCCESS" : @"FAILURE"));
	[pass release];
	[user release];
	[myRecord release];
	[myNode release];
	return retval;
}
*/
import "C"

import "os/user"

type darwinPasswordChecker struct{}

func NewPasswordChecker() PasswordChecker {
	return &darwinPasswordChecker{}
}

func (p *darwinPasswordChecker) CurrentUser() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

func (p *darwinPasswordChecker) CheckPassword(username, password string) bool {
	res := C.checkPassword(C.CString(username), C.CString(password))
	return bool(res)
}
