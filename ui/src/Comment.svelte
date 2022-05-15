<script>
	import CreateComment from "./CreateComment.svelte";
    import UpdateComment from "./UpdateComment.svelte";
    import {decodeJwtResponse, getAccessToken} from "./google.js";

    export let comment;

    let isReplying = false;

    function toggleReply() {
        isReplying = !isReplying
    }

    function isAuthor(email) {
        if (sessionStorage.hasOwnProperty("accessToken")) {
            const accessToken = getAccessToken();
            let responsePayload = decodeJwtResponse(accessToken);
            return responsePayload.email === email;
        } else {
            return false
        }
    }

    let isUpdating = false;

    function toggleUpdate() {
        isUpdating = !isUpdating
    }
</script>

{#if isUpdating}
	<UpdateComment id="{comment.id}" content="{comment.content}"/>
{:else}
	<div id="comment{comment.id}">
		<img src="{comment.profile_picture}" width="48" height="48">
		<strong>{comment.name}</strong> . {comment.created_at}<br>
		{comment.content}<br>
		<button type="button" class="btn btn-outline-dark replybtn" id="{comment.id}" on:click={toggleReply}>Reply</button>
	{#if isAuthor(comment.email)}
		<button type="button" class="btn btn-outline-dark replybtn" on:click={toggleUpdate}>Edit</button>
	{/if}
	</div>
{/if}

{#if isReplying}
	<CreateComment parentId={comment.id} innerHTML="Reply"/>
{/if}

{#if comment.comments}
	<ul>
		{#each comment.comments as c}
			<li>
				<svelte:self comment={c}/>
			</li>
		{/each}
	</ul>
{/if}