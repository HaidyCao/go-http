package windows

/*

#include <stdio.h>

extern void Test(const char* url, const char* user, const char* pwd);

void test_fubc(const char* url, const char* user, const char* pwd) {
	Text(url, user, pwd);
}

*/

import "C"

import "github.com/HaidyCao/go-http/http"

func Test(u *C.char, user *C.char, pwd *C.char) {
	url := C.GoString(u)
	username := C.GoString(user)
	password := C.GoString(pwd)
	client := &http.GoClient{
		Url:    url,
		Method: "POST",
		Transport: &http.GoHttpTransport{
			Username: username,
			Password: password,
			Ntlm:     true,
		},
		ContentType: "text/xml; charset=utf-8",
		PostData:    []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:m=\"http://schemas.microsoft.com/exchange/services/2006/messages\" xmlns:t=\"http://schemas.microsoft.com/exchange/services/2006/types\"><soap:Header><t:RequestServerVersion Version=\"Exchange2010_SP2\"/></soap:Header><soap:Body><m:FindItem Traversal=\"Shallow\"><m:ItemShape><t:BaseShape>IdOnly</t:BaseShape></m:ItemShape><m:ParentFolderIds><t:DistinguishedFolderId Id=\"calendar\"/></m:ParentFolderIds></m:FindItem></soap:Body></soap:Envelope>"),
	}

	client.AddHeaderNameAndValue("Accept", "text/xml")

	resp, _ := http.Request(client)

	str, _ := resp.GetBody().GetData()
	println(string(str))
}
