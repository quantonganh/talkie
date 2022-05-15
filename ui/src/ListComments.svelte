<script>
	import { comments, getCommentsByPost } from "./stores.js";
    import Comment from "./Comment.svelte";

    getCommentsByPost();
</script>

{#each $comments as comment}
	<Comment comment={comment}/>
{/each}
