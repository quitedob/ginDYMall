package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"douyin/pkg/utils/upload"    // Corrected path to pkg/utils/upload
	"douyin/pkg/utils/response"
	"douyin/middleware"
	"douyin/mylog" // Assuming your logger, adjust if it's from conf.GlobalConfig.Logger or similar
)

// UploadController handles file uploads.
type UploadController struct {
	OSSClient *upload.Client // Corrected type to upload.Client
}

// RegisterRoutes mounts upload routes.
func (uc *UploadController) RegisterRoutes(r *gin.Engine) {
	// Apply AuthMiddleware to the upload group
	uploadGroup := r.Group("/api/v1/upload", middleware.AuthMiddleware())
	{
		uploadGroup.POST("/file", uc.UploadFile)
		uploadGroup.POST("/files", uc.UploadFiles)
	}
}

// UploadFile handles single file uploads.
// @Summary      Upload a single file
// @Description  Uploads a single file to configured OSS (e.g., S3)
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to upload"
// @Param        path formData string false "Optional sub-path in bucket (e.g., 'avatars', 'products'). Defaults to 'general'."
// @Success      200 {object} response.APIResponse{data=object{url=string}} "Upload successful, returns file URL"
// @Failure      400 {object} response.APIResponse "Bad request (e.g., no file, invalid path)"
// @Failure      401 {object} response.APIResponse "Unauthorized"
// @Failure      500 {object} response.APIResponse "Internal server error (e.g., upload failed)"
// @Router       /upload/file [post]
func (uc *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "未选择文件或无效的文件字段 (No file chosen or invalid file field): "+err.Error())
		return
	}

	uploadPath := c.PostForm("path")
	if uploadPath == "" {
		uploadPath = "general" // Default path
	}

	if uc.OSSClient == nil {
		mylog.Error("OSSClient is not initialized in UploadController") // Ensure mylog is correctly initialized and available
		response.Fail(c, http.StatusInternalServerError, "文件服务未初始化 (File service not initialized)")
		return
	}

	url, err := uc.OSSClient.Upload(uploadPath, file)
	if err != nil {
		mylog.Errorf("Failed to upload file: %v", err)
		response.Fail(c, http.StatusInternalServerError, "上传文件失败 (Failed to upload file): "+err.Error())
		return
	}
	response.Success(c, gin.H{"url": url})
}

// UploadFiles handles multiple file uploads.
// @Summary      Upload multiple files
// @Description  Uploads multiple files to configured OSS (e.g., S3)
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        files formData []file true "Files to upload"
// @Param        path  formData string false "Optional sub-path in bucket for all files. Defaults to 'general'."
// @Success      200 {object} response.APIResponse{data=object{urls=[]string}} "Upload successful, returns file URLs"
// @Failure      400 {object} response.APIResponse "Bad request (e.g., no files)"
// @Failure      401 {object} response.APIResponse "Unauthorized"
// @Failure      500 {object} response.APIResponse "Internal server error (e.g., one or more uploads failed)"
// @Router       /upload/files [post]
func (uc *UploadController) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "获取上传表单失败 (Failed to get multipart form): "+err.Error())
		return
	}
	files := form.File["files"] // "files" is the field name for multiple files

	if len(files) == 0 {
		response.Fail(c, http.StatusBadRequest, "未选择任何文件 (No files chosen)")
		return
	}

	uploadPath := c.PostForm("path")
	if uploadPath == "" {
		uploadPath = "general" // Default path
	}

	if uc.OSSClient == nil {
		mylog.Error("OSSClient is not initialized in UploadController")
		response.Fail(c, http.StatusInternalServerError, "文件服务未初始化 (File service not initialized)")
		return
	}

	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := uc.OSSClient.Upload(uploadPath, file)
		if err != nil {
			mylog.Errorf("Failed to upload one of the files (%s): %v", file.Filename, err)
			response.Fail(c, http.StatusInternalServerError, fmt.Sprintf("上传文件 %s 失败 (Failed to upload file %s): %s", file.Filename, file.Filename, err.Error()))
			return
		}
		urls = append(urls, url)
	}
	response.Success(c, gin.H{"urls": urls})
}
