export function handleCredentialResponse(response) {
    authGoogle(response.credential)
        .then(token => {
            setAccessToken(token.access_token)
            setRefreshToken(token.refresh_token)
        })

    // decodeJwtResponse() is a custom function defined by you
    // to decode the credential response.
    let responsePayload = decodeJwtResponse(response.credential);
    console.log(responsePayload)

    console.log("ID: " + responsePayload.sub);
    console.log('Full Name: ' + responsePayload.name);
    console.log('Given Name: ' + responsePayload.given_name);
    console.log('Family Name: ' + responsePayload.family_name);
    console.log("Image URL: " + responsePayload.picture);
    console.log("Email: " + responsePayload.email);

    document.getElementById("commentContent").disabled = false;
    document.getElementById("commentBtn").disabled = false;
}

export function decodeJwtResponse(token) {
    let base64Url = token.split('.')[1]
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload)
}

async function authGoogle(credential) {
    let body = JSON.stringify({
        "credential": credential,
    })
    const response = await fetch('TALKIE_URL' + "/auth/google", {
        method: "post",
        headers: { "Content-Type": "application/json" },
        body: body
    });
    return response.json()
}

export function getAccessToken() {
    return sessionStorage.getItem("accessToken")
}

export function setAccessToken(token) {
    sessionStorage.setItem("accessToken", token)
}

export function getRefreshToken() {
    return sessionStorage.getItem("refreshToken")
}

export function setRefreshToken(token) {
    sessionStorage.setItem("refreshToken", token)
}
