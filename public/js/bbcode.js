// [tag]Affected text[/tag] (b, i, u...)
function wrapSimpleBBCode(e, before, after){
    var textarea = document.getElementById("article-text");
    var tag = e.getAttribute("id");
    var len = textarea.value.length;
    var start = textarea.selectionStart;
    var end = textarea.selectionEnd;
    var selection = textarea.value.substring(start, end);
    if (typeof before == "string") {
        selection = before+selection;
    }
    if (typeof after == "string") {
        selection = selection+after;
    }
    var replace = "["+tag+"]"+selection+"[/"+tag+"]";
    textarea.value = textarea.value.substring(0,start) + replace +
        textarea.value.substring(end,len);
}
// [tag]property[/tag] (img, youtube)
function wrapLinkBBCode(tag, link){
    var textarea = document.getElementById("article-text");
    if (tag == "youtube") {
        var found = link.match(/v=(.*)/);
        link = found && found[1] || link;
    }
    var len = textarea.value.length;
    var start = textarea.selectionStart;
    var replace = "["+tag+"]"+link.toString()+"[/"+tag+"]";
    textarea.value = textarea.value.substring(0,start) + replace +
        textarea.value.substring(start,len);
}
// Not implemented, for [tag=property]Affected text[/tag] style tags (font color, url...)
function wrapPropertyBBCode(tag, property){
    var textarea = document.getElementById("article-text");
    var tag = e.getAttribute("id");
    var len = textarea.value.length;
    var start = textarea.selectionStart;
    var end = textarea.selectionEnd;
    var selection = textarea.value.substring(start, end);
    var replace = "["+tag+"="+property+"]"+selection+"[/"+tag+"]";
    textarea.value = textarea.value.substring(0,start) + replace +
        textarea.value.substring(end,len);
}