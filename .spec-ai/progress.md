STATUS: COMPLETE

## Tasks

- [x] TASK 1: Add ProfilePictureURL field to User struct in domain/user.go
- [x] TASK 2: Create SQL migration for profile_picture_url column
- [x] TASK 3: Update repository queries (SaveUser, GetUserByEmail, FindUserByUuid, UpdateUser, ListUsers) to include profile_picture_url
- [x] TASK 4: Create upload_profile_picture.go usecase
- [x] TASK 5: Add UploadProfilePicture handler method to user_handler.go
- [x] TASK 6: Add error mappings for upload errors in error_mapper.go
- [x] TASK 7: Wire usecase, update handler constructor, register route and static file serving in main.go
- [x] TASK 8: Build and verify


## Compliance — Missing Requirements (fix these)
- [x] Requirement: File Validation - Added Content-Type validation (image/jpeg, image/png) in upload_profile_picture.go via validateContentType method
- [x] Requirement: API Contract & Error Handling - Added 'code' field (INVALID_FILE, USER_NOT_FOUND, INTERNAL_ERROR) to upload error responses via mapUploadErrorCode in user_handler.go
