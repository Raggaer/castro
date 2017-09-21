function froala(selector) {
    $(selector).froalaEditor({
        imageUploadURL: "/subtopic/admin/froala/image",
        imageUploadParams: { _csrf: csrfToken },
        imageAllowedTypes: ["png"],
        fileUploadURL: "/subtopic/admin/froala/file",
        fileUploadParams: { _csrf: csrfToken }
    });
}