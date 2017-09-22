function froala(selector) {
    $(selector).froalaEditor({
        imageManagerLoadURL: "/subtopic/admin/froala/manager/list",
        imageUploadURL: "/subtopic/admin/froala/upload/image",
        imageUploadParams: { _csrf: csrfToken },
        imageAllowedTypes: ["png"],
        fileUploadURL: "/subtopic/admin/froala/upload/file",
        fileUploadParams: { _csrf: csrfToken },
        imageManagerDeleteMethod: "POST",
        imageManagerDeleteURL: "/subtopic/admin/froala/manager/delete",
        imageManagerDeleteParams: { _csrf: csrfToken }
    });
}