<script>
	import {getCommentsByPost} from "./stores.js";
    import {getAccessToken} from "./google.js";

	export let id;
    export let content;

    export async function updateComment(id, content) {
        const accessToken = getAccessToken();

        let body = JSON.stringify({
            "content": content,
        })
        await fetch('TALKIE_URL' + "/comments/" + id, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
	            "Authorization": `Bearer ` + accessToken,
            },
            body: body
        })
            .then(resp => {
                if (resp.status === 200) {
                    return resp.json()
                } else {
                    console.log("Status: " + resp.status)
                }
            });

        await getCommentsByPost();
    }
</script>

<div class="form-group">
	<textarea class="form-control" id="commentContent" rows="3" bind:value={content}></textarea>
</div>
<div class="d-flex flex-row-reverse">
	<button type="button" class="btn btn-success" id="commentBtn" on:click={updateComment(id, content)}>Update</button>
</div>
