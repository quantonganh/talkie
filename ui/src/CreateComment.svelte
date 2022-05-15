<script>
	import {getCommentsByPost} from "./stores.js";
    import {getAccessToken} from "./google.js";

    export function isLoggedIn() {
        return sessionStorage.hasOwnProperty("accessToken")
    }

    export let shown = true;
	export let parentId;
    export let innerHTML;
    let content;
    let commentContent;

    async function createComment(parentId, content) {
        const accessToken = getAccessToken();

        let body = JSON.stringify({
            "post_slug": location.pathname.slice(1),
            "parent_id": parentId,
            "content": content,
        })
        await fetch('TALKIE_URL' + "/comments", {
            method: "POST",
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

        if (parentId === undefined) {
	        document.getElementById("commentContent").value = '';
        } else {
            shown = !shown
        }
    }
</script>

{#if shown}
	<div class="form-group">
		<textarea disabled={!isLoggedIn()} class="form-control" id="commentContent" rows="3" bind:this={commentContent} bind:value={content}></textarea>
	</div>
	<div class="d-flex flex-row-reverse">
		<button disabled={!isLoggedIn()} type="button" class="btn btn-success" id="commentBtn" on:click={createComment(parentId, content)}>{innerHTML}</button>
	</div>
{/if}